package utils

import "golang.org/x/crypto/bcrypt"

// CryptHash generates a bcrypt hash of the password combined with salt.
// It returns the hashed password and any error encountered during hashing.
func CryptHash(raw, salt string) (string, error) {
	combined := raw + salt
	hash, err := bcrypt.GenerateFromPassword([]byte(combined), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CryptHashCompare compares a raw password with salt against a bcrypt hash.
// It returns true if the password matches the hash, false otherwise.
func CryptHashCompare(raw, salt, hash string) bool {
	combined := raw + salt
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(combined))
	return err == nil
}
