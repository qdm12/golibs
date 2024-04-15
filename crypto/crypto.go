package crypto

import (
	"github.com/qdm12/golibs/crypto/random"
)

type Crypto struct {
	random *random.Random
}

func NewCrypto() *Crypto {
	return &Crypto{
		random: random.NewRandom(),
	}
}
