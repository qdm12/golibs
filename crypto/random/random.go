package random

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

const (
	lowercase     = "abcdefghijklmnopqrstuvwxyz"
	uppercase     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alpha         = lowercase + uppercase
	num           = "0123456789"
	alphaNum      = alpha + num
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	maxInt63      = 9223372036854775807
)

// Random has methods to generate random values
//go:generate mockgen -destination=mock_random/random.go . Random
type Random interface {
	GenerateRandomBytes(n int) ([]byte, error)
	GenerateRandomInt63() int64
	GenerateRandomInt(n int) int
	GenerateRandomAlphaNum(length uint64) string
	GenerateRandomNum(n uint64) string
}

// random implements Random
type random struct {
	randReader func(b []byte) error
}

// NewRandom returns a new Random object
func NewRandom() Random {
	return &random{
		randReader: randReader,
	}
}

func randReader(b []byte) error {
	l := len(b)
	n, err := rand.Read(b)
	if l != n {
		return fmt.Errorf("read %d random bytes instead of expected %d bytes", n, l)
	} else if err != nil {
		return err
	}
	return nil
}

// GenerateRandomBytes generates n random bytes
func (r *random) GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if err := r.randReader(b); err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateRandomInt63 returns a random 63 bit positive integer
func (r *random) GenerateRandomInt63() int64 {
	const int63Length = 32 // 256 bits
	b, err := r.GenerateRandomBytes(int63Length)
	if err != nil {
		panic(err.Error())
	}
	v := int64(binary.BigEndian.Uint64(b))
	if v < 0 {
		v = -v
	}
	return v
}

// GenerateRandomInt returns a random integer between 0 and n
func (r *random) GenerateRandomInt(n int) int {
	if n == 0 {
		return 0
	}
	return int(r.GenerateRandomInt63()) % n
}

// GenerateRandomAlphaNum returns a string of random alphanumeric characters of a specified length
func (r *random) GenerateRandomAlphaNum(length uint64) string {
	if length >= maxInt63 {
		panic("length argument cannot be bigger than 2^63 - 1")
	}
	n := int64(length)
	b := make([]byte, n)
	// GenerateRandomInt63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, r.GenerateRandomInt63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = r.GenerateRandomInt63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(alphaNum) {
			b[i] = alphaNum[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// GenerateRandomNum returns a string of random numeric characters of a specified length
func (r *random) GenerateRandomNum(n uint64) string {
	if n >= maxInt63 {
		panic("length argument cannot be bigger than 2^63 - 1")
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = num[int(r.GenerateRandomInt63())%len(num)]
	}
	return string(b)
}
