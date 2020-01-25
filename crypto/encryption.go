package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

// EncryptAES256 encrypts some plaintext with a key using AES and returns the ciphertext
func (c *crypto) EncryptAES256(plaintext []byte, key [32]byte) (ciphertext []byte, err error) {
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

// DecryptAES256 decrypts some ciphertext with a key using AES and returns the plaintext
func (c *crypto) DecryptAES256(ciphertext []byte, key [32]byte) (plaintext []byte, err error) {
	block, _ := aes.NewCipher(key[:])
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("DecryptAES: cipher size %d should be bigger than block size %d", len(ciphertext), aes.BlockSize)
	}
	iv := ciphertext[:aes.BlockSize]
	plaintext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintext, plaintext)
	return plaintext, nil
}
