package maphash

import (
	"hash/maphash"
	"math/rand"
)

var _ rand.Source = new(Source)

func New() *Source {
	return &Source{}
}

type Source struct{}

func (s *Source) Int63() int64 {
	v := new(maphash.Hash).Sum64()
	return int64(v >> 1)
}

func (s *Source) Seed(_ int64) {}
