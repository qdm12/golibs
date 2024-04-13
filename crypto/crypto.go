package crypto

import (
	"github.com/qdm12/golibs/crypto/random"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/sha3"
)

// Crypto contains methods to run cryptographic operations.
type Crypto interface {
	// EncryptAES256 uses the AES algorithm with a 256 bits key to encrypt a plaintext and returns a ciphertext
	EncryptAES256(plaintext []byte, key [32]byte) (ciphertext []byte, err error)
	// DecryptAES256 uses the AES algorithm with a 256 bits key to decrypt a ciphertext and returns a plaintext
	DecryptAES256(ciphertext []byte, key [32]byte) (plaintext []byte, err error)
	// ShakeSum256 uses the SHA3 Shake Hash 256 algorithm to produce a 512 bits digest from some data
	ShakeSum256(data []byte) (digest [shakeSum256DigestSize]byte, err error)
	// Argon2ID uses the Argon2ID hash algorithm to produce a 512 bits digest from some data
	Argon2ID(data []byte, time, memory uint32) (digest [argon2IDDigestSize]byte)
	// Checksumize adds a checksum to some data using the SHA2 256 algorithm
	Checksumize(data []byte) (checksumedData []byte, err error)
	// Dehecksumize verifies the checksum matches the data and removes it if so, it uses the SHA2 256 algorithm
	Dechecksumize(checksumData []byte) (data []byte, err error)
	// SignEd25519 signs a message with a 512 bits signing key and returns the signature
	SignEd25519(message []byte, signingKey [signKeySize]byte) (signature []byte)
	// NewSalt generates a random 256 bits salt
	NewSalt() (salt [saltSize]byte, err error)
}

type crypto struct {
	shakeHashFactory func() sha3.ShakeHash
	argon2ID         func(password []byte, salt []byte, time uint32, memory uint32, threads uint8, keyLen uint32) []byte
	random           random.Randomer
}

func NewCrypto() Crypto {
	return &crypto{
		shakeHashFactory: func() sha3.ShakeHash { return sha3.NewShake256() }, //nolint:gocritic
		argon2ID:         argon2.IDKey,
		random:           random.NewRandom(),
	}
}
