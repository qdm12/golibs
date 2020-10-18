package network

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/golibs/crypto/random/mock_random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewClient(t *testing.T) {
	t.Parallel()
	var c Client
	require.NotPanics(t, func() {
		c = NewClient(time.Second)
	})
	assert.NotNil(t, c)
	_, ok := c.(*client)
	assert.True(t, ok)
}

func Test_Close(t *testing.T) {
	t.Parallel()
	client := NewClient(time.Nanosecond)
	assert.NotPanics(t, func() {
		client.Close()
	})
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (rtf roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return rtf(r)
}

func newMockRoundTripper(t *testing.T, expectedRequest *http.Request, response *http.Response, err error) http.RoundTripper {
	return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, expectedRequest.URL, r.URL)
		return response, err
	})
}

type errReader struct{}

func (e *errReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("read error")
}

func (e *errReader) Close() error { return nil }

func Test_Do(t *testing.T) {
	t.Parallel()
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()
	tests := map[string]struct {
		ctx       context.Context
		request   *http.Request
		response  *http.Response
		clientErr error
		content   []byte
		status    int
		err       error
	}{
		"no error": {
			ctx:     context.Background(),
			request: &http.Request{Method: http.MethodGet, URL: &url.URL{}},
			response: &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBufferString("body")),
				StatusCode: http.StatusOK,
			},
			content: []byte("body"),
			status:  http.StatusOK,
		},
		"http error": {
			ctx:       context.Background(),
			request:   &http.Request{Method: http.MethodGet, URL: &url.URL{}},
			clientErr: fmt.Errorf("http error"),
			err:       fmt.Errorf(`Get "": http error`),
		},
		"context canceled": {
			ctx:     canceledCtx,
			request: &http.Request{Method: http.MethodGet, URL: &url.URL{}},
			err:     fmt.Errorf("context canceled"),
		},
		"body read error": {
			ctx:     context.Background(),
			request: &http.Request{Method: http.MethodGet, URL: &url.URL{}},
			response: &http.Response{
				Body:       &errReader{},
				StatusCode: http.StatusOK,
			},
			status: http.StatusOK,
			err:    fmt.Errorf("read error"),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			httpClient := &http.Client{
				Transport: newMockRoundTripper(t, tc.request, tc.response, tc.clientErr),
			}
			c := &client{
				httpClient: httpClient,
			}
			content, status, err := c.Do(tc.ctx, tc.request)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.status, status)
			assert.Equal(t, tc.content, content)
		})
	}
}

func Test_UseRandomUserAgent(t *testing.T) {
	t.Parallel()
	setter := UseRandomUserAgent()
	options := getOptions{}
	setter(&options)
	assert.Equal(t, getOptions{randomUserAgent: true}, options)
}

func Test_GetContent(t *testing.T) {
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()
	t.Parallel()
	tests := map[string]struct {
		ctx             context.Context
		URL             string
		setters         []GetSetter
		expectedRequest *http.Request
		response        *http.Response
		clientErr       error
		content         []byte
		status          int
		err             error
	}{
		"bad url": {
			ctx: context.Background(),
			URL: "\n",
			err: fmt.Errorf(`parse "\n": net/url: invalid control character in URL`),
		},
		"http error": {
			ctx: context.Background(),
			URL: "https://domain.com",
			expectedRequest: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Scheme: "https",
					Host:   "domain.com",
				},
			},
			clientErr: fmt.Errorf("error"),
			err:       fmt.Errorf(`Get "https://domain.com": error`),
		},
		"context canceled": {
			ctx: canceledCtx,
			URL: "https://domain.com",
			expectedRequest: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Scheme: "https",
					Host:   "domain.com",
				},
			},
			err: fmt.Errorf("context canceled"),
		},
		"no error": {
			ctx: context.Background(),
			URL: "https://domain.com",
			expectedRequest: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Scheme: "https",
					Host:   "domain.com",
				},
			},
			response: &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBufferString("body")),
				StatusCode: http.StatusOK,
			},
			content: []byte("body"),
			status:  http.StatusOK,
		},
		"no error with random user agent": {
			ctx:     context.Background(),
			URL:     "https://domain.com",
			setters: []GetSetter{UseRandomUserAgent()},
			expectedRequest: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Scheme: "https",
					Host:   "domain.com",
				},
				Header: http.Header{"User-Agent": []string{"b"}},
			},
			response: &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBufferString("body")),
				StatusCode: http.StatusOK,
			},
			content: []byte("body"),
			status:  http.StatusOK,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			random := mock_random.NewMockRandom(ctrl)
			httpClient := &http.Client{
				Transport: newMockRoundTripper(t, tc.expectedRequest, tc.response, tc.clientErr),
			}
			userAgents := []string{"a", "b", "c"}
			random.EXPECT().GenerateRandomInt(len(userAgents)).Return(1)
			c := &client{
				httpClient: httpClient,
				userAgents: userAgents,
				random:     random,
			}
			content, status, err := c.Get(tc.ctx, tc.URL, tc.setters...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.status, status)
			assert.Equal(t, tc.content, content)
		})
	}
}
