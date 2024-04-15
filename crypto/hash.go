package crypto

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/sha3"
)

var (
	ErrBytesWrittenUnexpected = errors.New("number of bytes written is unexpected")
	ErrBytesReadUnexpected    = errors.New("number of bytes read is unexpected")
)

func (c *Crypto) ShakeSum256(data []byte) (digest [shakeSum256DigestSize]byte, err error) {
	buf := make([]byte, shakeSum256DigestSize)
	shakeHash := sha3.NewShake256()
	n, err := shakeHash.Write(data)
	if n != len(data) {
		return digest, fmt.Errorf("%w: %d bytes written instead of %d",
			ErrBytesWrittenUnexpected, n, len(data))
	} else if err != nil {
		return digest, err
	}
	n, err = shakeHash.Read(buf)
	if n != shakeSum256DigestSize {
		return digest, fmt.Errorf("%w: %d bytes read instead of %d",
			ErrBytesReadUnexpected, n, shakeSum256DigestSize)
	} else if err != nil {
		return digest, err
	}
	copy(digest[:], buf)
	return digest, nil
}

func (c *Crypto) Argon2ID(data []byte, time, memory uint32) (digest [argon2IDDigestSize]byte) {
	buf := argon2.IDKey(data, nil, time, memory, 1, argon2IDDigestSize)
	copy(digest[:], buf)
	return digest
}
