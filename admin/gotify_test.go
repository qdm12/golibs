package admin

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewGotify(t *testing.T) {
	t.Parallel()
	g := NewGotify(url.URL{}, "abc", &http.Client{})
	assert.NotNil(t, g)
}

// func Test_Ping(t *testing.T) {
// 	t.Parallel()
// 	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("2"))
// 	}))
// 	defer s.Close()
// 	URL, err := url.Parse("http://" + s.Listener.Addr().String())
// 	require.NoError(t, err)
// 	httpClient := &http.Client{
// 		Transport: &http.Transport{
// 			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
// 				return net.Dial(network, s.Listener.Addr().String())
// 			},
// 		},
// 	}
// 	g := NewGotify(*URL, "a", httpClient)
// 	assert.NoError(t, g.Ping())
// }
