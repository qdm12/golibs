package syncpool

import (
	"crypto/rand"
	"encoding/binary"
	"sync"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Rand

// Rand implements a fast and thread safe pseudo-random number generator
// seeded with the crypto/rand reader. It uses some memory to store multiple
// instances of pseudo random generator states to scale with CPU cores.
type Rand interface {
	// Uint32 returns a pseudo random uint32 number.
	Uint32() uint32
	// Uint32n returns a pseudo random uint32 number in the range [0..maxN).
	Uint32n(maxN uint32) uint32
}

type syncPoolRand struct {
	pool *sync.Pool
}

func New() Rand {
	return &syncPoolRand{
		pool: &sync.Pool{
			New: func() interface{} {
				return &prng{
					n: randomUint32(),
				}
			},
		},
	}
}

func (s *syncPoolRand) Uint32() uint32 {
	v := s.pool.Get()
	r := v.(*prng)
	n := r.uint32()
	s.pool.Put(r)
	return n
}

func (s *syncPoolRand) Uint32n(maxN uint32) uint32 {
	n := s.Uint32()
	// See http://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/
	return uint32((uint64(n) * uint64(maxN)) >> 32) //nolint:gomnd
}

type prng struct {
	n uint32 // seeded with randomUint32
}

//nolint:gomnd
func (r *prng) uint32() uint32 {
	// See https://en.wikipedia.org/wiki/Xorshift
	r.n ^= r.n << 13
	r.n ^= r.n >> 17
	r.n ^= r.n << 5
	return r.n
}

func randomUint32() uint32 {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return binary.BigEndian.Uint32(b)
}
