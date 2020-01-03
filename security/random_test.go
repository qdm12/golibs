package security

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateRandomInt63(t *testing.T) {
	t.Parallel()
	for i := 0; i < 20; i++ {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			out := GenerateRandomInt63()
			assert.GreaterOrEqual(t, out, int64(0))
		})
	}
}

func Test_GenerateRandomInt(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		n int
	}{
		"generates random int modulo 0":  {0},
		"generates random int modulo 1":  {1},
		"generates random int modulo 3":  {3},
		"generates random int modulo 50": {50},
	}
	for name, tc := range tests {
		for i := 0; i < 20; i++ {
			t.Run(fmt.Sprintf("%s %d", name, i), func(t *testing.T) {
				tc := tc
				t.Parallel()
				out := GenerateRandomInt(tc.n)
				assert.GreaterOrEqual(t, out, 0)
				assert.LessOrEqual(t, out, tc.n)
			})
		}

	}
}

func Test_GenerateRandomAlphaNum(t *testing.T) {
	t.Parallel()
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
			s := GenerateRandomAlphaNum(tc.n)
			assert.Regexp(t, tc.regex, s)
		})
	}
	t.Run("panics with length argument too large", func(t *testing.T) {
		assert.PanicsWithValue(t, "length argument cannot be bigger than 2^63 - 1", func() { GenerateRandomAlphaNum(9223372036854775808) })
	})
}

func Test_GenerateRandomNum(t *testing.T) {
	t.Parallel()
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
			s := GenerateRandomNum(tc.n)
			assert.Regexp(t, tc.regex, s)
		})
	}
	t.Run("panics with length argument too large", func(t *testing.T) {
		assert.PanicsWithValue(t, "length argument cannot be bigger than 2^63 - 1", func() { GenerateRandomNum(9223372036854775808) })
	})
}
