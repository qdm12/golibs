package crypto

import (
	"bytes"
	"fmt"
)

// Checksumize adds a Shake 256 checksum to some data
func (c *crypto) Checksumize(data []byte) (checksumedData []byte, err error) {
	digest, err := c.ShakeSum256(data)
	if err != nil {
		return nil, fmt.Errorf("Checksumize: %w", err)
	}
	checksum := digest[:checksumLength]
	return append(data, checksum...), nil
}

// Dechecksumize verifies the Shake 256 checksum of some data
func (c *crypto) Dechecksumize(checksumData []byte) (data []byte, err error) {
	L := len(checksumData)
	if L < checksumLength {
		return nil, fmt.Errorf("checksumed data %v not long enough to contain the checksum", checksumData)
	}
	checksum := checksumData[L-4:]
	data = checksumData[:L-4]
	digest, err := c.ShakeSum256(data)
	if err != nil {
		return nil, fmt.Errorf("Dechecksumize: %w", err)
	}
	checksum2 := digest[:4]
	if !bytes.Equal(checksum, checksum2) {
		return nil, fmt.Errorf("checksum verification failed (%v and %v)", checksum, checksum2)
	}
	return data, nil
}
