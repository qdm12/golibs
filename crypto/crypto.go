package crypto

import (
	"github.com/qdm12/golibs/crypto/random"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/sha3"
)

type Crypto struct {
	shakeHashFactory func() sha3.ShakeHash
	argon2ID         func(password []byte, salt []byte, time uint32, memory uint32, threads uint8, keyLen uint32) []byte
	random           *random.Random
}

func NewCrypto() *Crypto {
	return &Crypto{
		shakeHashFactory: func() sha3.ShakeHash { return sha3.NewShake256() }, //nolint:gocritic
		argon2ID:         argon2.IDKey,
		random:           random.NewRandom(),
	}
}
