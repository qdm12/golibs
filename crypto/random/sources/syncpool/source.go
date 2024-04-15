package syncpool

import (
	"fmt"
	"math/rand"
	"sync"
)

var _ rand.Source = new(Source)

func New() *Source {
	return &Source{
		pool: &sync.Pool{
			New: func() interface{} {
				return newGenerator()
			},
		},
	}
}

type Source struct {
	pool *sync.Pool
}

func (s *Source) Int63() int64 {
	generator, ok := s.pool.Get().(*generator)
	if !ok {
		panic(fmt.Sprintf("unexpected type %T", generator))
	}

	v := generator.uint64()

	s.pool.Put(generator)

	return int64(v >> 1)
}

func (s *Source) Seed(_ int64) {}
