package security

import (
	mathrand "math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generateRandomAlphaNum(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		source mathrand.Source
		n      uint64
		s      string
	}{
		"generates string of length 0": {
			mathrand.NewSource(0),
			0,
			"",
		},
		"generates string of length 1": {
			mathrand.NewSource(0),
			1,
			"b",
		},
		"generates string of length 10": {
			mathrand.NewSource(0),
			10,
			"haJ8lRczqb",
		},
		"generates string of length 10 with another source": {
			mathrand.NewSource(1),
			10,
			"p1LGIehp1s",
		},
		"generates string of length 50": {
			mathrand.NewSource(0),
			50,
			"P0IV8ngpT1cY1kB0VsRO7QvIp2xLpBSHZ9Bbbl00haJ8lRczqb",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc := tc
			t.Parallel()
			s := generateRandomAlphaNum(tc.n, tc.source)
			assert.Equal(t, tc.s, s)
		})
	}
	t.Run("panics with length argument too large", func(t *testing.T) {
		assert.PanicsWithValue(t, "length argument cannot be bigger than 2^63 - 1", func() { GenerateRandomAlphaNum(9223372036854775808) })
	})
}

func Test_generateRandomNum(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		source mathrand.Source
		n      uint64
		s      string
	}{
		"generates string of length 0": {
			mathrand.NewSource(0),
			0,
			"",
		},
		"generates string of length 1": {
			mathrand.NewSource(0),
			1,
			"5",
		},
		"generates string of length 10": {
			mathrand.NewSource(0),
			10,
			"5274233696",
		},
		"generates string of length 10 with another source": {
			mathrand.NewSource(1),
			10,
			"0111708869",
		},
		"generates string of length 50": {
			mathrand.NewSource(0),
			50,
			"52742336966898095198669060670904277435924082803451",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc := tc
			t.Parallel()
			s := generateRandomNum(tc.n, tc.source)
			assert.Equal(t, tc.s, s)
		})
	}
	t.Run("panics with length argument too large", func(t *testing.T) {
		assert.PanicsWithValue(t, "length argument cannot be bigger than 2^63 - 1", func() { GenerateRandomNum(9223372036854775808) })
	})
}

func Test_generateRandomInt(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		source mathrand.Source
		n      int
		out    int
	}{
		"generates random int modulo 0": {
			mathrand.NewSource(0),
			0,
			0,
		},
		"generates random int modulo 1": {
			mathrand.NewSource(0),
			1,
			0,
		},
		"generates random int modulo 3": {
			mathrand.NewSource(0),
			3,
			2,
		},
		"generates random int modulo 10 with another source": {
			mathrand.NewSource(1),
			10,
			0,
		},
		"generates random int modulo 50": {
			mathrand.NewSource(0),
			50,
			5,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc := tc
			t.Parallel()
			out := generateRandomInt(tc.n, tc.source)
			assert.Equal(t, tc.out, out)
		})
	}
}
