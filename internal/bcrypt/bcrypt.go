package bcrypt

import "golang.org/x/crypto/bcrypt"

const SALT_BCRYPT = 12

// Hash generates a bcrypt hash from the given password.
// It takes a byte slice as input and returns the hashed password as a string.
// If an error occurs during the hashing process, it returns an empty string and the error.
func Hash(p []byte) (pass string, err error) {
	b, err := bcrypt.GenerateFromPassword(p, SALT_BCRYPT)

	if err != nil {
		return
	}

	pass = string(b)

	return
}

// Compare compares a hashed password with a plaintext password.
// It uses bcrypt.CompareHashAndPassword to perform the comparison.
// If the passwords match, it returns nil. Otherwise, it returns an error.
func Compare(hash []byte, pass []byte) (err error) {
	err = bcrypt.CompareHashAndPassword(hash, pass)

	return err
}
