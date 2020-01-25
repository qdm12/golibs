package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_EncryptAES(t *testing.T) {
	t.Parallel()
	c := NewCrypto()
	tests := map[string]struct {
		plaintext    []byte
		key          [32]byte
		err          error
		cipherLength int
	}{
		"works with data": {
			[]byte("The quick brown fox jumps over the lazy dog"),
			[32]byte{12, 32, 77, 57, 96, 15, 221, 211, 241, 242, 12, 168, 0, 126, 145, 199, 208, 41, 59, 28, 195, 145, 10, 59, 248, 178, 230, 29, 160, 242, 107, 202},
			nil,
			59,
		},
		"works with short data": {
			[]byte{100},
			[32]byte{12, 32, 77, 57, 96, 15, 221, 211, 241, 242, 12, 168, 0, 126, 145, 199, 208, 41, 59, 28, 195, 145, 10, 59, 248, 178, 230, 29, 160, 242, 107, 202},
			nil,
			17,
		},
		// TODO mock
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ciphertext, err := c.EncryptAES256(tc.plaintext, tc.key)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Len(t, ciphertext, tc.cipherLength)
		})
	}
}

func Test_DecryptAES(t *testing.T) {
	t.Parallel()
	c := NewCrypto()
	tests := map[string]struct {
		ciphertext []byte
		key        [32]byte
		plaintext  []byte
		err        error
	}{
		"works with encrypted data": {
			[]byte{83, 170, 176, 118, 26, 89, 244, 96, 153, 247, 56, 128, 34, 168, 187, 43, 194, 158, 217, 64, 91, 46, 91, 227, 110, 43, 228, 145, 23, 119, 223, 24, 154, 224, 157, 27, 97, 219, 135, 142, 226, 132, 103, 33, 31, 48, 117, 232, 216, 20, 169, 106, 169, 209, 101, 42, 43, 10, 222},
			[32]byte{12, 32, 77, 57, 96, 15, 221, 211, 241, 242, 12, 168, 0, 126, 145, 199, 208, 41, 59, 28, 195, 145, 10, 59, 248, 178, 230, 29, 160, 242, 107, 202},
			[]byte("The quick brown fox jumps over the lazy dog"),
			nil,
		},
		"works with short encrypted data": {
			[]byte{46, 142, 130, 63, 245, 220, 21, 167, 184, 40, 28, 130, 135, 236, 73, 36, 229},
			[32]byte{12, 32, 77, 57, 96, 15, 221, 211, 241, 242, 12, 168, 0, 126, 145, 199, 208, 41, 59, 28, 195, 145, 10, 59, 248, 178, 230, 29, 160, 242, 107, 202},
			[]byte{100},
			nil,
		},
		"empty data": {
			[]byte{},
			[32]byte{12, 32, 77, 57, 96, 15, 221, 211, 241, 242, 12, 168, 0, 126, 145, 199, 208, 41, 59, 28, 195, 145, 10, 59, 248, 178, 230, 29, 160, 242, 107, 202},
			nil,
			fmt.Errorf("DecryptAES: cipher size 0 should be bigger than block size 16"),
		},
		"data too short": {
			[]byte{45, 156, 61},
			[32]byte{12, 32, 77, 57, 96, 15, 221, 211, 241, 242, 12, 168, 0, 126, 145, 199, 208, 41, 59, 28, 195, 145, 10, 59, 248, 178, 230, 29, 160, 242, 107, 202},
			nil,
			fmt.Errorf("DecryptAES: cipher size 3 should be bigger than block size 16"),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out, err := c.DecryptAES256(tc.ciphertext, tc.key)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.plaintext, out)
		})
	}
}

func Test_EncryptDecryptAES256(t *testing.T) {
	t.Parallel()
	c := NewCrypto()
	tests := map[string]struct {
		plaintext  []byte
		key        [32]byte
		encryptErr error
		decryptErr error
	}{
		"works with encrypted data": {
			[]byte("The quick brown fox jumps over the lazy dog"),
			[32]byte{12, 32, 77, 57, 96, 15, 221, 211, 241, 242, 12, 168, 0, 126, 145, 199, 208, 41, 59, 28, 195, 145, 10, 59, 248, 178, 230, 29, 160, 242, 107, 202},
			nil,
			nil,
		},
		"works with short data": {
			[]byte("Short"),
			[32]byte{12, 32, 77, 57, 96, 15, 221, 211, 241, 242, 12, 168, 0, 126, 145, 199, 208, 41, 59, 28, 195, 145, 10, 59, 248, 178, 230, 29, 160, 242, 107, 202},
			nil,
			nil,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ciphertext, err := c.EncryptAES256(tc.plaintext, tc.key)
			if tc.encryptErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.encryptErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			if err != nil {
				return
			}
			plaintext, err := c.DecryptAES256(ciphertext, tc.key)
			if tc.decryptErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.decryptErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			if err == nil {
				assert.Equal(t, tc.plaintext, plaintext)
			}
		})
	}
}
