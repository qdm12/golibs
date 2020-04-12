package crypto

import (
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/sha3"

	"github.com/qdm12/golibs/crypto/random"
)

//go:generate mockgen -destination=mock_crypto/crypto.go . Crypto
type Crypto interface {
	EncryptAES256(plaintext []byte, key [32]byte) (ciphertext []byte, err error)
	DecryptAES256(ciphertext []byte, key [32]byte) (plaintext []byte, err error)
	ShakeSum256(data []byte) (digest [shakeSum256DigestSize]byte, err error)
	Argon2ID(data []byte, time, memory uint32) (digest [argon2IDDigestSize]byte)
	Checksumize(data []byte) (checksumedData []byte, err error)
	Dechecksumize(checksumData []byte) (data []byte, err error)
	SignEd25519(message []byte, signingKey [signKeySize]byte) (signature []byte)
	NewSalt() (salt [saltSize]byte, err error)
}

type crypto struct {
	shakeHashFactory func() sha3.ShakeHash
	argon2ID         func(password []byte, salt []byte, time uint32, memory uint32, threads uint8, keyLen uint32) []byte
	random           random.Random
}

func NewCrypto() Crypto {
	return &crypto{
		shakeHashFactory: func() sha3.ShakeHash { return sha3.NewShake256() }, //nolint:gocritic
		argon2ID:         argon2.IDKey,
		random:           random.NewRandom(),
	}
}
