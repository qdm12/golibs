package security

import (
	"crypto/rand"
	"encoding/binary"
)

const lowercase = "abcdefghijklmnopqrstuvwxyz"
const uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const alpha = lowercase + uppercase
const num = "0123456789"
const alphaNum = alpha + num
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateRandomInt63 returns a random 63 bit positive integer
func GenerateRandomInt63() int64 {
	b, _ := GenerateRandomBytes(32) // 256 bits
	v := int64(binary.BigEndian.Uint64(b))
	if v < 0 {
		v = -v
	}
	return v
}

// GenerateRandomInt returns a random integer between 0 and n
func GenerateRandomInt(n int) int {
	if n == 0 {
		return 0
	}
	r := int(GenerateRandomInt63())
	if r < 0 {
		r = -r
	}
	return int(r % n)
}

// GenerateRandomAlphaNum returns a string of random alphanumeric characters of a specified length
func GenerateRandomAlphaNum(length uint64) string {
	if length >= 9223372036854775807 {
		panic("length argument cannot be bigger than 2^63 - 1")
	}
	n := int64(length)
	b := make([]byte, n)
	// GenerateRandomInt63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, GenerateRandomInt63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = GenerateRandomInt63(), letterIdxMax
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
func GenerateRandomNum(n uint64) string {
	if n >= 9223372036854775807 {
		panic("length argument cannot be bigger than 2^63 - 1")
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = num[int(GenerateRandomInt63())%len(num)]
	}
	return string(b)
}
