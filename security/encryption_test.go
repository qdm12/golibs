package security

import (
	"errors"
	"fmt"
	"testing"

	"github.com/qdm12/golibs/helpers"
	"github.com/stretchr/testify/assert"
)

func Test_EncryptAES(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		plaintext    []byte
		key          []byte
		err          error
		cipherLength int
	}{
		"works with data": {
			[]byte("The quick brown fox jumps over the lazy dog"),
			[]byte("gKeg4J7wPChk8DQJC9U86rJMEMSalIog"),
			nil,
			59,
		},
		"works with short data": {
			[]byte{100},
			[]byte("gKeg4J7wPChk8DQJC9U86rJMEMSalIog"),
			nil,
			17,
		},
		"key of wrong size": {
			[]byte("The quick brown fox jumps over the lazy dog"),
			[]byte{10, 24, 5},
			errors.New("EncryptAES: crypto/aes: invalid key size 3"),
			0,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out, err := EncryptAES(tc.plaintext, tc.key)
			helpers.AssertErrorsEqual(t, tc.err, err)
			assert.Len(t, out, tc.cipherLength)
		})
	}
}

func Test_DecryptAES(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		ciphertext []byte
		key        []byte
		plaintext  []byte
		err        error
	}{
		"works with encrypted data": {
			[]byte{45, 156, 61, 33, 8, 244, 174, 30, 156, 30, 131, 75, 217, 22, 121, 189, 72, 62, 82, 154, 0, 49, 186, 217, 115, 16, 186, 206, 209, 69, 82, 225, 207, 10, 111, 223, 10, 26, 249, 5, 123, 124, 124, 43, 198, 158, 213, 227, 208, 21, 119, 52, 73, 187, 137, 15, 255, 166, 73},
			[]byte("gKeg4J7wPChk8DQJC9U86rJMEMSalIog"),
			[]byte("The quick brown fox jumps over the lazy dog"),
			nil,
		},
		"empty data": {
			[]byte{},
			[]byte("gKeg4J7wPChk8DQJC9U86rJMEMSalIog"),
			nil,
			fmt.Errorf("DecryptAES: cipher size 0 should be bigger than block size 16"),
		},
		"data too short": {
			[]byte{45, 156, 61},
			[]byte("gKeg4J7wPChk8DQJC9U86rJMEMSalIog"),
			nil,
			fmt.Errorf("DecryptAES: cipher size 3 should be bigger than block size 16"),
		},
		"key of wrong size": {
			[]byte{45, 156, 61, 33, 8, 244, 174, 30, 156, 30, 131, 75, 217, 22, 121, 189, 72, 62, 82, 154, 0, 49, 186, 217, 115, 16, 186, 206, 209, 69, 82, 225, 207, 10, 111, 223, 10, 26, 249, 5, 123, 124, 124, 43, 198, 158, 213, 227, 208, 21, 119, 52, 73, 187, 137, 15, 255, 166, 73},
			[]byte{10, 24, 5},
			nil,
			fmt.Errorf("DecryptAES: crypto/aes: invalid key size 3"),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out, err := DecryptAES(tc.ciphertext, tc.key)
			helpers.AssertErrorsEqual(t, tc.err, err)
			assert.Equal(t, tc.plaintext, out)
		})
	}
}

func Test_EncryptDecryptAES(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		plaintext  []byte
		key        []byte
		encryptErr error
		decryptErr error
	}{
		"works with encrypted data": {
			[]byte("The quick brown fox jumps over the lazy dog"),
			[]byte("gKeg4J7wPChk8DQJC9U86rJMEMSalIog"),
			nil,
			nil,
		},
		"works with short data": {
			[]byte("Short"),
			[]byte("gKeg4J7wPChk8DQJC9U86rJMEMSalIog"),
			nil,
			nil,
		},
		"key of wrong size": {
			[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 50, 86, 180, 114, 48, 74, 17, 217, 5, 12, 186, 135, 216, 204, 243, 201, 149, 15, 247, 104, 87, 30, 106, 111, 229, 232, 100, 77, 147, 233, 134, 159, 237, 198, 101, 54, 41, 22, 95, 18, 64, 128, 36},
			[]byte{10, 24, 5},
			fmt.Errorf("EncryptAES: crypto/aes: invalid key size 3"),
			nil,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ciphertext, err := EncryptAES(tc.plaintext, tc.key)
			if !helpers.AssertErrorsEqual(t, tc.encryptErr, err) {
				t.FailNow()
			}
			if err != nil {
				return
			}
			plaintext, err := DecryptAES(ciphertext, tc.key)
			helpers.AssertErrorsEqual(t, tc.decryptErr, err)
			if err == nil {
				assert.Equal(t, tc.plaintext, plaintext)
			}
		})
	}
}
