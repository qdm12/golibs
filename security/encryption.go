package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// EncryptAES encrypts some plaintext with a key using AES and returns the ciphertext
func EncryptAES(plaintext []byte, key []byte) (ciphertext []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("EncryptAES: %w", err)
	}
	ciphertext = make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, fmt.Errorf("EncryptAES: %w", err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext, nil
}

// DecryptAES decrypts some ciphertext with a key using AES and returns the plaintext
func DecryptAES(ciphertext []byte, key []byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("DecryptAES: %w", err)
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("DecryptAES: cipher size %d should be bigger than block size %d", len(ciphertext), aes.BlockSize)
	}
	iv := ciphertext[:aes.BlockSize]
	plaintext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintext, plaintext)
	return plaintext, nil
}
