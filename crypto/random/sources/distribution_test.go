package sources

import (
	"testing"

	"github.com/qdm12/golibs/crypto/random/sources/maphash"
	"github.com/qdm12/golibs/crypto/random/sources/syncpool"
)

func Test_Maphash_Distribution(t *testing.T) {
	t.Parallel()

	source := maphash.New()
	f := func(n int) (out int) {
		out = int(source.Int63())
		if out < 0 {
			out = -out
		}
		return out % n
	}

	isUniformlyDistributed(t, f)
}

func Test_SyncPool_Distribution(t *testing.T) {
	t.Parallel()

	source := syncpool.New()
	f := func(n int) (out int) {
		out = int(source.Int63())
		if out < 0 {
			out = -out
		}
		return out % n
	}

	isUniformlyDistributed(t, f)
}

func isUniformlyDistributed(t *testing.T, f func(n int) int) {
	t.Helper()

	const iterations = 100000
	const maxValue = 30

	numberToCount := make(map[int]int, maxValue)

	for i := 0; i < iterations; i++ {
		out := f(maxValue)
		numberToCount[out]++
	}

	targetGenerationsPerNumber := iterations / maxValue
	const maxPercentDivergence = 0.07

	for number := 0; number < maxValue; number++ {
		count := numberToCount[number]

		diff := targetGenerationsPerNumber - count
		if diff < 0 {
			diff = -diff
		}
		divergence := float64(diff) / float64(targetGenerationsPerNumber)
		if divergence > maxPercentDivergence {
			t.Errorf("Number %d was generated %d times, with a %.2f%% divergence from the target %d",
				number, count, 100*divergence, targetGenerationsPerNumber)
		}
	}
}
