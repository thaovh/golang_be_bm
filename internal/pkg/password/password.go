package password

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost is the default bcrypt cost
	DefaultCost = 12
)

// Hash generates a bcrypt hash from the password
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Verify compares a password with a hash
func Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// MustHash generates a hash and panics on error (use with caution)
func MustHash(password string) string {
	hash, err := Hash(password)
	if err != nil {
		panic(err)
	}
	return hash
}

