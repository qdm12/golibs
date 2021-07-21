package crypto

import (
	"bytes"
	"errors"
	"fmt"
)

// Checksumize adds a Shake 256 checksum to some data.
func (c *crypto) Checksumize(data []byte) (checksumedData []byte, err error) {
	digest, err := c.ShakeSum256(data)
	if err != nil {
		return nil, err
	}
	checksum := digest[:checksumLength]
	return append(data, checksum...), nil
}

var (
	ErrChecksumLength   = errors.New("checksum is not long enough")
	ErrChecksumMismatch = errors.New("checkum mismatch")
)

// Dechecksumize verifies the Shake 256 checksum of some data.
func (c *crypto) Dechecksumize(checksumData []byte) (data []byte, err error) {
	L := len(checksumData)
	if L < checksumLength {
		return nil, fmt.Errorf("%w: expected at least %d bytes and got only %d: %v",
			ErrChecksumLength, checksumLength, L, checksumData)
	}
	checksum := checksumData[L-4:]
	data = checksumData[:L-4]
	digest, err := c.ShakeSum256(data)
	if err != nil {
		return nil, err
	}
	checksum2 := digest[:4]
	if !bytes.Equal(checksum, checksum2) {
		return nil, fmt.Errorf("%w: expected %v but got %v",
			ErrChecksumMismatch, checksum, checksum2)
	}
	return data, nil
}
