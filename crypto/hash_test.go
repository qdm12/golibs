package crypto

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/qdm12/golibs/crypto/mocks"
)

func Test_Shake256(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		data   []byte
		digest [shakeSum256DigestSize]byte
		err    error
	}{
		"no data": {
			nil,
			[shakeSum256DigestSize]byte{0x46, 0xb9, 0xdd, 0x2b, 0xb, 0xa8, 0x8d, 0x13, 0x23, 0x3b, 0x3f, 0xeb, 0x74, 0x3e, 0xeb, 0x24, 0x3f, 0xcd, 0x52, 0xea, 0x62, 0xb8, 0x1b, 0x82, 0xb5, 0xc, 0x27, 0x64, 0x6e, 0xd5, 0x76, 0x2f, 0xd7, 0x5d, 0xc4, 0xdd, 0xd8, 0xc0, 0xf2, 0x0, 0xcb, 0x5, 0x1, 0x9d, 0x67, 0xb5, 0x92, 0xf6, 0xfc, 0x82, 0x1c, 0x49, 0x47, 0x9a, 0xb4, 0x86, 0x40, 0x29, 0x2e, 0xac, 0xb3, 0xb7, 0xc4, 0xbe},
			nil,
		},
		"empty data": {
			[]byte{},
			[shakeSum256DigestSize]byte{0x46, 0xb9, 0xdd, 0x2b, 0xb, 0xa8, 0x8d, 0x13, 0x23, 0x3b, 0x3f, 0xeb, 0x74, 0x3e, 0xeb, 0x24, 0x3f, 0xcd, 0x52, 0xea, 0x62, 0xb8, 0x1b, 0x82, 0xb5, 0xc, 0x27, 0x64, 0x6e, 0xd5, 0x76, 0x2f, 0xd7, 0x5d, 0xc4, 0xdd, 0xd8, 0xc0, 0xf2, 0x0, 0xcb, 0x5, 0x1, 0x9d, 0x67, 0xb5, 0x92, 0xf6, 0xfc, 0x82, 0x1c, 0x49, 0x47, 0x9a, 0xb4, 0x86, 0x40, 0x29, 0x2e, 0xac, 0xb3, 0xb7, 0xc4, 0xbe},
			nil,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			c := NewCrypto()
			digest, err := c.ShakeSum256(tc.data)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.digest, digest)
		})
	}
}

func Test_Shake256_Mocked(t *testing.T) {
	t.Parallel()
	type mockShakehashWrite struct {
		call bool
		n    int
		err  error
	}
	type mockShakehashRead struct {
		call bool
		n    int
		err  error
	}
	tests := map[string]struct {
		shakehashWrite mockShakehashWrite
		shakehashRead  mockShakehashRead
		data           []byte
		digest         [shakeSum256DigestSize]byte
		err            error
	}{
		"shakehash wrong written number": {
			mockShakehashWrite{call: true, n: 5},
			mockShakehashRead{},
			[]byte("data"),
			[shakeSum256DigestSize]byte{},
			fmt.Errorf("Shake256: 5 bytes written instead of 4"),
		},
		"shakehash write error": {
			mockShakehashWrite{call: true, n: 4, err: fmt.Errorf("error")},
			mockShakehashRead{},
			[]byte("data"),
			[shakeSum256DigestSize]byte{},
			fmt.Errorf("Shake256: error"),
		},
		"shakehash wrong read number": {
			mockShakehashWrite{call: true, n: 4},
			mockShakehashRead{call: true, n: 4},
			[]byte("data"),
			[shakeSum256DigestSize]byte{},
			fmt.Errorf("Shake256: 4 bytes read instead of 64"),
		},
		"shakehash read error": {
			mockShakehashWrite{call: true, n: 4},
			mockShakehashRead{call: true, n: 64, err: fmt.Errorf("error")},
			[]byte("data"),
			[shakeSum256DigestSize]byte{},
			fmt.Errorf("Shake256: error"),
		},
		"shakehash success": {
			mockShakehashWrite{call: true, n: 4},
			mockShakehashRead{call: true, n: 64},
			[]byte("data"),
			[shakeSum256DigestSize]byte{},
			nil,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockShakeHash := &mocks.ShakeHash{}
			mockShakeHash.On("Reset").Once()
			if tc.shakehashWrite.call {
				mockShakeHash.On("Write", tc.data).
					Return(tc.shakehashWrite.n, tc.shakehashWrite.err).Once()
			}
			if tc.shakehashRead.call {
				mockShakeHash.On("Read", make([]byte, shakeSum256DigestSize)).
					Return(tc.shakehashRead.n, tc.shakehashRead.err).Once()
			}
			c := crypto{shakeHash: mockShakeHash}
			digest, err := c.ShakeSum256(tc.data)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.digest, digest)
			mockShakeHash.AssertExpectations(t)
		})
	}
}
