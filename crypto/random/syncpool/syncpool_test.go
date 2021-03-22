package syncpool

import (
	"testing"
)

func Test_syncPoolRand_Uint32(t *testing.T) {
	rand := New()

	itemsCounts := make(map[uint32]int)
	const tries = 1e4 // this has to be large enough
	for i := 0; i < tries; i++ {
		n := rand.Uint32()
		itemsCounts[n]++
	}

	uniquenessRatio := 100 * float64(len(itemsCounts)) / tries
	t.Log(uniquenessRatio)
	if uniquenessRatio < 95 {
		t.Errorf("only %.2f%% of numbers are unique", uniquenessRatio)
	}

	const maxRepetition = 3
	for number, count := range itemsCounts {
		if count > maxRepetition {
			t.Errorf("number %d was generated %d times", number, count)
		}
	}
}

func Test_syncPoolRand_Uint32n(t *testing.T) {
	rand := New()

	const nMax = 20
	const tries = 1e5 // this has to be large enough
	const maxSkew = 5 // percents

	itemsCounts := make(map[uint32]int)
	for i := 0; i < tries; i++ {
		n := rand.Uint32n(nMax)
		itemsCounts[n]++
	}

	avg := tries / nMax
	for number, count := range itemsCounts {
		p := 100 * (float64(count) - avg) / avg
		if p < 0 {
			p = -p
		}
		if p > maxSkew {
			t.Errorf("skew is too high for number %d: %.2f%%", number, p)
		}
	}
}
