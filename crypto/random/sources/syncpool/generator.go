package syncpool

import (
	"crypto/rand"
	"encoding/binary"
)

func newGenerator() *generator {
	return &generator{
		n: makeSeed(),
	}
}

type generator struct {
	n uint64
}

func makeSeed() (seed uint64) {
	b := make([]byte, 8) //nolint:gomnd
	_, _ = rand.Read(b)
	return binary.BigEndian.Uint64(b)
}

func (g *generator) uint64() uint64 {
	g.n ^= g.n << 13 //nolint:gomnd
	g.n ^= g.n >> 7  //nolint:gomnd
	g.n ^= g.n << 17 //nolint:gomnd
	return g.n
}
