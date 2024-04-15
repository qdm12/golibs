package crypto

import (
	"github.com/qdm12/golibs/crypto/random"
	"golang.org/x/crypto/argon2"
)

type Crypto struct {
	argon2ID func(password []byte, salt []byte, time uint32, memory uint32, threads uint8, keyLen uint32) []byte
	random   *random.Random
}

func NewCrypto() *Crypto {
	return &Crypto{
		argon2ID: argon2.IDKey,
		random:   random.NewRandom(),
	}
}
