package security

import (
	mathrand "math/rand"
	"time"
)

var randomSource = mathrand.NewSource(time.Now().UnixNano())

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

// GenerateRandomAlphaNum returns a string of random alphanumeric characters of a specified length
func GenerateRandomAlphaNum(length uint64) string {
	return generateRandomAlphaNum(length, randomSource)
}

// GenerateRandomNum returns a string of random numeric characters of a specified length
func GenerateRandomNum(length uint64) string {
	return generateRandomAlphaNum(length, randomSource)
}

// GenerateRandomInt returns a random integer between 0 and n
func GenerateRandomInt(n int) int {
	return generateRandomInt(n, randomSource)
}

func generateRandomAlphaNum(length uint64, source mathrand.Source) string {
	if length >= 9223372036854775807 {
		panic("length argument cannot be bigger than 2^63 - 1")
	}
	n := int64(length)
	b := make([]byte, n)
	// A source.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, source.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = source.Int63(), letterIdxMax
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

func generateRandomNum(n uint64, source mathrand.Source) string {
	if n >= 9223372036854775807 {
		panic("length argument cannot be bigger than 2^63 - 1")
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = num[int(source.Int63())%len(num)]
	}
	return string(b)
}

func generateRandomInt(n int, source mathrand.Source) int {
	if n == 0 {
		return 0
	}
	return int(source.Int63() % int64(n))
}
