package crypto

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Chechsumize(t *testing.T) {
	t.Parallel()
	c := NewCrypto()
	tests := map[string]struct {
		data         []byte
		checksumData []byte
		err          error
	}{
		"some data": {
			[]byte{215, 168, 251, 179, 7, 215, 128, 148},
			[]byte{0xd7, 0xa8, 0xfb, 0xb3, 0x7, 0xd7, 0x80, 0x94, 0x3b, 0x16, 0x90, 0x87},
			nil,
		},
		"empty data": {
			[]byte{},
			[]byte{0x46, 0xb9, 0xdd, 0x2b},
			nil,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			checksumData, err := c.Checksumize(tc.data)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.checksumData, checksumData)
		})
	}
}

func Test_Dechecksumize(t *testing.T) {
	t.Parallel()
	c := NewCrypto()
	tests := map[string]struct {
		checksumData []byte
		data         []byte
		err          error
	}{
		"some data checksummed": {
			[]byte{0xd7, 0xa8, 0xfb, 0xb3, 0x7, 0xd7, 0x80, 0x94, 0x3b, 0x16, 0x90, 0x87},
			[]byte{215, 168, 251, 179, 7, 215, 128, 148},
			nil,
		},
		"empty data checksummed": {
			[]byte{0x46, 0xb9, 0xdd, 0x2b},
			[]byte{},
			nil,
		},
		"data with bad checksum": {
			[]byte{0xe7, 0xa8, 0xfb, 0xb3, 0x7, 0xd7, 0x80, 0x94, 0x3b, 0x16, 0x90, 0x87},
			nil,
			errors.New("checkum mismatch: expected [59 22 144 135] but got [95 138 74 156]"),
		},
		"data not long enough": {
			[]byte{0xe7, 0xe7, 0xe7},
			nil,
			errors.New("checksum is not long enough: expected at least 4 bytes and got only 3: [231 231 231]"),
		},
		"empty data": {
			[]byte{},
			nil,
			errors.New("checksum is not long enough: expected at least 4 bytes and got only 0: []"),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			data, err := c.Dechecksumize(tc.checksumData)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.data, data)
		})
	}
}
