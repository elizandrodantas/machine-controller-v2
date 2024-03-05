package security

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/elizandrodantas/machine-controller-v2/internal/util"
)

type crypto struct {
	Key []byte
}

// Crypto is a function that creates a new instance of the crypto struct.
// If the key parameter is nil, it generates a random 16-byte key.
// The generated crypto struct contains the key used for encryption and decryption.
func Crypto(key []byte) *crypto {
	if key == nil {
		key = util.RandomByte(16)
	}

	return &crypto{key}
}

// Encrypt encrypts the given plaintext using AES encryption in CTR mode.
// It takes the plaintext and initialization vector (IV) as input and returns the ciphertext.
// If an error occurs during encryption, it returns nil and the corresponding error.
func (y *crypto) Encrypt(plaintext, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(y.Key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)

	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)

	return ciphertext, nil
}

// Decrypt decrypts the given ciphertext using the provided initialization vector (IV).
// It returns the decrypted plaintext or an error if decryption fails.
func (y *crypto) Decrypt(ciphertext, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(y.Key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)

	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}
