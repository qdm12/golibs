package security

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

// between 1 and 32
const checksumLength = 4

func init() {
	if checksumLength < 1 || checksumLength > 32 {
		panic(fmt.Sprintf("Checksum length %d must be between 1 and 32 bytes", checksumLength))
	}
}

// Checksumize adds a SHA256 checksum to some data
func Checksumize(data []byte) (checksumData []byte) {
	digest := sha256.Sum256(data)
	checksum := digest[:checksumLength]
	checksumData = append(data, checksum...)
	return checksumData
}

// Dechecksumize verifies the SHA256 checksum of some data
func Dechecksumize(checksumData []byte) (data []byte, err error) {
	L := len(checksumData)
	if L < 4 {
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
