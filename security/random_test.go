package security

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GenerateRandomBytes(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		n          int
		randReader func(b []byte) error
		err        error
	}{
		"no error": {
			10,
			func(b []byte) error {
				_, err := rand.Read(b)
				return err
			},
			nil,
		},
		"error": {
			0,
			func(b []byte) error {
				return fmt.Errorf("error")
			},
			fmt.Errorf("error"),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r := &RandomImpl{
				randReader: tc.randReader,
			}
			out, err := r.GenerateRandomBytes(tc.n)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
			assert.Len(t, out, tc.n)
		})
	}
}

func Test_GenerateRandomInt63(t *testing.T) {
	t.Parallel()
	r := NewRandom()
	var previousValue int64
	for i := 0; i < 50; i++ {
		out := r.GenerateRandomInt63()
		assert.GreaterOrEqual(t, out, int64(0))
		assert.NotEqual(t, out, previousValue)
		previousValue = out
	}
	t.Run("panics from rand.Read error", func(t *testing.T) {
		t.Parallel()
		r := &RandomImpl{
			randReader: func(b []byte) error {
				return fmt.Errorf("error")
			},
		}
		assert.PanicsWithValue(t, "error", func() { r.GenerateRandomInt63() })
	})
}

func Test_GenerateRandomInt(t *testing.T) {
	t.Parallel()
	r := NewRandom()
	tests := map[string]struct {
		n int
	}{
		"generates random int modulo 0":  {0},
		"generates random int modulo 1":  {1},
		"generates random int modulo 3":  {3},
		"generates random int modulo 50": {50},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			for i := 0; i < 50; i++ {
				out := r.GenerateRandomInt(tc.n)
				assert.GreaterOrEqual(t, out, 0)
				assert.LessOrEqual(t, out, tc.n)
			}
		})
	}
}

func Test_GenerateRandomAlphaNum(t *testing.T) {
	t.Parallel()
	r := NewRandom()
	tests := map[string]struct {
		n     uint64
		regex string
	}{
		"generates string of length 0": {
			0,
			`^$`,
		},
		"generates string of length 1": {
			1,
			`^[a-zA-Z0-9]$`,
		},
		"generates string of length 10": {
			10,
			"^[a-zA-Z0-9]{10}$",
		},
		"generates string of length 50": {
			50,
			"^[a-zA-Z0-9]{50}$",
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := r.GenerateRandomAlphaNum(tc.n)
			assert.Regexp(t, tc.regex, s)
		})
	}
	t.Run("panics with length argument too large", func(t *testing.T) {
		t.Parallel()
		assert.PanicsWithValue(t, "length argument cannot be bigger than 2^63 - 1", func() { r.GenerateRandomAlphaNum(maxInt63 + 1) })
	})
}

func Test_GenerateRandomNum(t *testing.T) {
	t.Parallel()
	r := NewRandom()
	tests := map[string]struct {
		n     uint64
		regex string
	}{
		"generates string of length 0": {
			0,
			"^$",
		},
		"generates string of length 1": {
			1,
			"^[0-9]$",
		},
		"generates string of length 10": {
			10,
			"^[0-9]{10}$",
		},
		"generates string of length 10 with another source": {
			10,
			"[0-9]{10}",
		},
		"generates string of length 50": {
			50,
			"[0-9]{50}",
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := r.GenerateRandomNum(tc.n)
			assert.Regexp(t, tc.regex, s)
		})
	}
	t.Run("panics with length argument too large", func(t *testing.T) {
		t.Parallel()
		assert.PanicsWithValue(t, "length argument cannot be bigger than 2^63 - 1", func() { r.GenerateRandomNum(maxInt63 + 1) })
	})
}
