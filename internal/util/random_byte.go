package util

import "crypto/rand"

// RandomByte generates a slice of random bytes with the specified length.
// It uses the crypto/rand package to generate secure random bytes.
// If an error occurs during the generation process, it returns a slice containing a single zero byte.
func RandomByte(len int) []byte {
	randomBytes := make([]byte, len)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return []byte{0}
	}

	return randomBytes
}
