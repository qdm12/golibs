package crypto

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"
)

func Test_Shake256(t *testing.T) {
	t.Parallel()
	//nolint:lll
	tests := map[string]struct {
		data   []byte
		digest [shakeSum256DigestSize]byte
		err    error
	}{
		"no data": {
			data:   nil,
			digest: [shakeSum256DigestSize]byte{0x46, 0xb9, 0xdd, 0x2b, 0xb, 0xa8, 0x8d, 0x13, 0x23, 0x3b, 0x3f, 0xeb, 0x74, 0x3e, 0xeb, 0x24, 0x3f, 0xcd, 0x52, 0xea, 0x62, 0xb8, 0x1b, 0x82, 0xb5, 0xc, 0x27, 0x64, 0x6e, 0xd5, 0x76, 0x2f, 0xd7, 0x5d, 0xc4, 0xdd, 0xd8, 0xc0, 0xf2, 0x0, 0xcb, 0x5, 0x1, 0x9d, 0x67, 0xb5, 0x92, 0xf6, 0xfc, 0x82, 0x1c, 0x49, 0x47, 0x9a, 0xb4, 0x86, 0x40, 0x29, 0x2e, 0xac, 0xb3, 0xb7, 0xc4, 0xbe},
			err:    nil,
		},
		"empty data": {
			data:   []byte{},
			digest: [shakeSum256DigestSize]byte{0x46, 0xb9, 0xdd, 0x2b, 0xb, 0xa8, 0x8d, 0x13, 0x23, 0x3b, 0x3f, 0xeb, 0x74, 0x3e, 0xeb, 0x24, 0x3f, 0xcd, 0x52, 0xea, 0x62, 0xb8, 0x1b, 0x82, 0xb5, 0xc, 0x27, 0x64, 0x6e, 0xd5, 0x76, 0x2f, 0xd7, 0x5d, 0xc4, 0xdd, 0xd8, 0xc0, 0xf2, 0x0, 0xcb, 0x5, 0x1, 0x9d, 0x67, 0xb5, 0x92, 0xf6, 0xfc, 0x82, 0x1c, 0x49, 0x47, 0x9a, 0xb4, 0x86, 0x40, 0x29, 0x2e, 0xac, 0xb3, 0xb7, 0xc4, 0xbe},
			err:    nil,
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

//go:generate mockgen -destination=mockShakeHash_test.go -package=crypto golang.org/x/crypto/sha3 ShakeHash
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
			fmt.Errorf("number of bytes written is unexpected: 5 bytes written instead of 4"),
		},
		"shakehash write error": {
			mockShakehashWrite{call: true, n: 4, err: fmt.Errorf("error")},
			mockShakehashRead{},
			[]byte("data"),
			[shakeSum256DigestSize]byte{},
			fmt.Errorf("error"),
		},
		"shakehash wrong read number": {
			mockShakehashWrite{call: true, n: 4},
			mockShakehashRead{call: true, n: 4},
			[]byte("data"),
			[shakeSum256DigestSize]byte{},
			fmt.Errorf("number of bytes read is unexpected: 4 bytes read instead of 64"),
		},
		"shakehash read error": {
			mockShakehashWrite{call: true, n: 4},
			mockShakehashRead{call: true, n: 64, err: fmt.Errorf("error")},
			[]byte("data"),
			[shakeSum256DigestSize]byte{},
			fmt.Errorf("error"),
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
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockShakeHash := NewMockShakeHash(mockCtrl)
			if tc.shakehashWrite.call {
				mockShakeHash.EXPECT().Write(tc.data).
					Return(tc.shakehashWrite.n, tc.shakehashWrite.err).Times(1)
			}
			if tc.shakehashRead.call {
				mockShakeHash.EXPECT().Read(make([]byte, shakeSum256DigestSize)).
					Return(tc.shakehashRead.n, tc.shakehashRead.err).Times(1)
			}
			c := crypto{shakeHashFactory: func() sha3.ShakeHash { return mockShakeHash }}
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

func Test_Argon2ID(t *testing.T) {
	t.Parallel()
	//nolint:lll
	tests := map[string]struct {
		data   []byte
		time   uint32
		memory uint32
		digest [argon2IDDigestSize]byte
	}{
		"no data": {
			data:   nil,
			time:   1,
			memory: 1,
			digest: [argon2IDDigestSize]byte{0xa3, 0x9d, 0x27, 0x7, 0x50, 0x31, 0x60, 0x36, 0xad, 0xc2, 0xe7, 0x7e, 0x6, 0x74, 0xb0, 0x92, 0x7d, 0x7a, 0x9b, 0xec, 0x56, 0x2e, 0xf5, 0x92, 0x3, 0xdd, 0xdb, 0xb9, 0x9f, 0xa1, 0xa7, 0xe3, 0xfa, 0x42, 0xde, 0xaf, 0x24, 0x26, 0xd3, 0x1b, 0x53, 0xa, 0xe5, 0x62, 0x50, 0xa2, 0x77, 0xa7, 0xac, 0x3a, 0x9a, 0x55, 0x6f, 0x1e, 0x92, 0x53, 0x11, 0x77, 0x3, 0xd6, 0x9a, 0xcd, 0x9e, 0x2e},
		},
		"no data with time 2, memory 1": {
			data:   nil,
			time:   2,
			memory: 1,
			digest: [argon2IDDigestSize]byte{0xf6, 0x4, 0xff, 0xaa, 0x20, 0x19, 0xec, 0x99, 0x4b, 0x3f, 0x51, 0x26, 0x46, 0x1a, 0xa, 0xcf, 0xb4, 0x11, 0x54, 0x14, 0x16, 0x2, 0x54, 0xc, 0x8d, 0xf7, 0xef, 0xb3, 0x7d, 0xd0, 0x14, 0x41, 0x39, 0xaa, 0xad, 0xa4, 0x7e, 0x3b, 0x4a, 0xf7, 0x62, 0x69, 0xe9, 0x1f, 0x54, 0xb9, 0x64, 0x83, 0xb7, 0x7c, 0x55, 0xc0, 0xc3, 0x28, 0x2b, 0xfe, 0x2f, 0x58, 0x9b, 0x5c, 0x80, 0x54, 0x97, 0xcd},
		},
		"data with time 2, memory 0": {
			data:   []byte{0, 1},
			time:   2,
			memory: 0,
			digest: [argon2IDDigestSize]byte{0x2f, 0x58, 0xfe, 0x52, 0x94, 0xd5, 0x41, 0x6c, 0x15, 0x5, 0x97, 0x91, 0x2a, 0x8f, 0xe9, 0x3b, 0x55, 0x27, 0xc2, 0x4a, 0x7e, 0xc0, 0xa7, 0xc2, 0x86, 0xd, 0x18, 0xe, 0x1b, 0xa0, 0xf7, 0xf1, 0x16, 0x1e, 0x4b, 0x58, 0xb5, 0xa7, 0xae, 0x76, 0xb1, 0x4, 0xa9, 0x4b, 0xe7, 0x93, 0x39, 0x84, 0xb3, 0xe0, 0x16, 0xd8, 0xc7, 0x96, 0x67, 0x3, 0xef, 0xd6, 0x97, 0xf6, 0x1d, 0x4e, 0xb5, 0x30},
		},
		"data with time 2, memory 2": {
			data:   []byte{0, 1},
			time:   2,
			memory: 2,
			digest: [argon2IDDigestSize]byte{0xd0, 0xe1, 0xea, 0x64, 0x25, 0x9f, 0x79, 0x72, 0x39, 0xea, 0x33, 0x2a, 0xa, 0xc6, 0x8b, 0xf4, 0x74, 0xe5, 0x45, 0xf2, 0x74, 0x6, 0xfb, 0x36, 0xb6, 0xf1, 0x2a, 0xc0, 0x5b, 0x5e, 0x1a, 0x8f, 0x2a, 0x7e, 0xb1, 0x41, 0x94, 0x60, 0x45, 0xcb, 0x50, 0xea, 0xf7, 0x62, 0xa1, 0x48, 0x28, 0x4f, 0xf6, 0xe1, 0x3b, 0x3b, 0x91, 0x4c, 0x93, 0xdf, 0x18, 0xa3, 0x1d, 0x38, 0xd8, 0xb1, 0x9, 0xe0},
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			c := NewCrypto()
			digest := c.Argon2ID(tc.data, tc.time, tc.memory)
			assert.Equal(t, tc.digest, digest)
		})
	}
}
