package crypto

import "crypto/ed25519"

func (c *crypto) SignEd25519(message []byte, signingKey [signKeySize]byte) (signature []byte) {
	privateKey := ed25519.PrivateKey(signingKey[:])
	return ed25519.Sign(privateKey, message)
}
