package random

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Randomer

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

// Randomer has methods to generate random values from the cryptographic
// randomness reader. All should be considered slow and thread safe.
type Randomer interface {
	// GenerateRandomBytes generates a byte slice of n random bytes
	GenerateRandomBytes(n int) ([]byte, error)
	// GenerateRandomInt63 generates a random 64 bits signed integer
	GenerateRandomInt63() int64
	// GenerateRandomInt generates a random signed integer
	GenerateRandomInt(n int) int
	// GenerateRandomAlphaNum generates a random alphanumeric string of n characters
	GenerateRandomAlphaNum(n uint64) string
	// GenerateRandomNum generates a random numeric string of n characters
	GenerateRandomNum(n uint64) string
}

// Random implements Randomer.
type Random struct {
	randReader func(b []byte) error
}

// NewRandom returns a new Random object.
func NewRandom() *Random {
	return &Random{
		randReader: randReader,
	}
}

var ErrRandReadBytesUnexpected = errors.New("read an unexpected number of random bytes")

func randReader(b []byte) error {
	n, err := rand.Read(b)
	if err != nil {
		return err
	} else if len(b) != n {
		return fmt.Errorf("%w: %d bytes instead of expected %d bytes",
			ErrRandReadBytesUnexpected, n, len(b))
	}
	return nil
}

// GenerateRandomBytes generates n random bytes.
func (r *Random) GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if err := r.randReader(b); err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateRandomInt63 returns a random 63 bit positive integer.
func (r *Random) GenerateRandomInt63() int64 {
	const int63Length = 8
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

// GenerateRandomInt returns a random integer between 0 and n.
func (r *Random) GenerateRandomInt(n int) (result int) {
	if n == 0 {
		return 0
	}
	result = int(r.GenerateRandomInt63()) % n
	if result < 0 {
		result = -result
	}
	return result
}

// GenerateRandomAlphaNum returns a string of random alphanumeric characters of a specified length.
func (r *Random) GenerateRandomAlphaNum(length uint64) string {
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

// GenerateRandomNum returns a string of random numeric characters of a specified length.
func (r *Random) GenerateRandomNum(n uint64) string {
	if n >= maxInt63 {
		panic("length argument cannot be bigger than 2^63 - 1")
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = num[int(r.GenerateRandomInt63())%len(num)]
	}
	return string(b)
}
