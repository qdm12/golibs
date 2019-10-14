package security

import (
	"errors"
	"testing"

	"github.com/qdm12/golibs/helpers"
	"github.com/stretchr/testify/assert"
)

func Test_Chechsumize(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		data         []byte
		checksumData []byte
	}{
		"some data": {
			[]byte{215, 168, 251, 179, 7, 215, 128, 148},
			[]byte{215, 168, 251, 179, 7, 215, 128, 148, 99, 19, 14, 134},
		},
		"empty data": {
			[]byte{},
			[]byte{227, 176, 196, 66},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc := tc
			t.Parallel()
			out := Checksumize(tc.data)
			assert.Equal(t, tc.checksumData, out)
		})
	}
}

func Test_Dechecksumize(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		checksumData []byte
		data         []byte
		err          error
	}{
		"some data checksummed": {
			[]byte{215, 168, 251, 179, 7, 215, 128, 148, 99, 19, 14, 134},
			[]byte{215, 168, 251, 179, 7, 215, 128, 148},
			nil,
		},
		"empty data checksummed": {
			[]byte{227, 176, 196, 66},
			[]byte{},
			nil,
		},
		"data with bad checksum": {
			[]byte{214, 168, 251, 179, 7, 215, 128, 148, 99, 19, 14, 134},
			nil,
			errors.New("checksum verification failed ([99 19 14 134] and [113 192 120 43])"),
		},
		"data not long enough": {
			[]byte{227, 176, 196},
			nil,
			errors.New("checksumed data [227 176 196] not long enough to contain the checksum"),
		},
		"empty data": {
			[]byte{},
			nil,
			errors.New("checksumed data [] not long enough to contain the checksum"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc := tc
			t.Parallel()
			out, err := Dechecksumize(tc.checksumData)
			helpers.AssertErrorsEqual(t, tc.err, err)
			assert.Equal(t, tc.data, out)
		})
	}
}
