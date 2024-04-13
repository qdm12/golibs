package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
)

// EncryptAES256 encrypts some plaintext with a key using AES and returns the ciphertext.
func (c *Crypto) EncryptAES256(plaintext []byte, key [32]byte) (ciphertext []byte, err error) {
	block, _ := aes.NewCipher(key[:])
	ciphertext = make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	randBytes, err := c.random.GenerateRandomBytes(len(iv))
	copy(iv, randBytes)
	if err != nil {
		return nil, fmt.Errorf("EncryptAES: %w", err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext, nil
}

var ErrCiphertextTooSmall = errors.New("ciphertext is too small")

// DecryptAES256 decrypts some ciphertext with a key using AES and returns the plaintext.
func (c *Crypto) DecryptAES256(ciphertext []byte, key [32]byte) (plaintext []byte, err error) {
	block, _ := aes.NewCipher(key[:])
	if len(ciphertext) < aes.BlockSize {
		return nil,
			fmt.Errorf("%w: is only %d bytes and must be at the %d bytes",
				ErrCiphertextTooSmall, len(ciphertext), aes.BlockSize)
	}
	iv := ciphertext[:aes.BlockSize]
	plaintext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintext, plaintext)
	return plaintext, nil
}
