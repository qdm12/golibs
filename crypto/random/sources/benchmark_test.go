package sources

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/qdm12/golibs/crypto/random/sources/maphash"
	"github.com/qdm12/golibs/crypto/random/sources/syncpool"
)

func Benchmark_MapHash(b *testing.B) {
	source := maphash.New()
	benchPerCoreConfigs(b, func(b *testing.B) { //nolint:thelper
		b.RunParallel(func(b *testing.PB) {
			for b.Next() {
				out := int(source.Int63())
				if out < 0 {
					out = -out
				}
				_ = out % 100
			}
		})
	})
}

func Benchmark_SyncPool(b *testing.B) {
	source := syncpool.New()
	benchPerCoreConfigs(b, func(b *testing.B) { //nolint:thelper
		b.RunParallel(func(b *testing.PB) {
			for b.Next() {
				out := int(source.Int63())
				if out < 0 {
					out = -out
				}
				_ = out % 100
			}
		})
	})
}

func benchPerCoreConfigs(b *testing.B, f func(b *testing.B)) {
	b.Helper()
	coreConfigs := []int{1, 2, 4, 8, 12, 18, 24}
	for _, n := range coreConfigs {
		n := n
		name := fmt.Sprintf("%d cores", n)
		b.Run(name, func(b *testing.B) {
			runtime.GOMAXPROCS(n)
			f(b)
		})
	}
}
