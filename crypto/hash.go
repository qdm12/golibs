package crypto

import (
	"fmt"
)

func (c *crypto) ShakeSum256(data []byte) (digest [shakeSum256DigestSize]byte, err error) {
	buf := make([]byte, shakeSum256DigestSize)
	c.shakeHash.Reset()
	n, err := c.shakeHash.Write(data)
	if n != len(data) {
		return digest, fmt.Errorf("Shake256: %d bytes written instead of %d", n, len(data))
	} else if err != nil {
		return digest, fmt.Errorf("Shake256: %w", err)
	}
	n, err = c.shakeHash.Read(buf)
	if n != shakeSum256DigestSize {
		return digest, fmt.Errorf("Shake256: %d bytes read instead of %d", n, shakeSum256DigestSize)
	} else if err != nil {
		return digest, fmt.Errorf("Shake256: %w", err)
	}
	copy(digest[:], buf)
	return digest, nil
}
