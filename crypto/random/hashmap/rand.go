package hashmap

import (
	"hash/maphash"
	"math/rand"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Rand

// Rand implements a fast and thread safe pseudo-random number generator
// seeded with the goroutine local storage.
type Rand interface {
	// Int returns a pseudo random int number.
	Int() int
	// Intn returns a pseudo random int number in the range [0..maxN).
	Intn(maxN int) int
}

type hashMapRand struct{}

func New() Rand {
	return &hashMapRand{}
}

func (h *hashMapRand) Int() int {
	return rand.New(source{}).Int() //nolint:gosec
}

func (h *hashMapRand) Intn(maxN int) int {
	return rand.New(source{}).Intn(maxN) //nolint:gosec
}

type source struct{}

func (source) Uint64() uint64 {
	return new(maphash.Hash).Sum64()
}

func (source) Int63() int64 {
	v := new(maphash.Hash).Sum64()
	return int64(v >> 1)
}

func (source) Seed(_ int64) {}
