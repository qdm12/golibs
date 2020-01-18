package security

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

// Checksumize adds a SHA256 checksum to some data
func Checksumize(data []byte) []byte {
	digest := sha256.Sum256(data)
	const checksumLength = 4 // between 1 and 32
	checksum := digest[:checksumLength]
	return append(data, checksum...)
}

// Dechecksumize verifies the SHA256 checksum of some data
func Dechecksumize(checksumData []byte) (data []byte, err error) {
	const minChecksumDataLength = 4
	L := len(checksumData)
	if L < minChecksumDataLength {
		return nil, fmt.Errorf("checksumed data %v not long enough to contain the checksum", checksumData)
	}
	checksum := checksumData[L-4:]
	data = checksumData[:L-4]
	digest := sha256.Sum256(data)
	checksum2 := digest[:4]
	if !bytes.Equal(checksum, checksum2) {
		return nil, fmt.Errorf("checksum verification failed (%v and %v)", checksum, checksum2)
	}
	return data, nil
}
