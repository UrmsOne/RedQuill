package common

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashPassword provides a simple SHA-256 based hashing. Replace with bcrypt/argon2 in production.
func HashPassword(plain string) string {
	if plain == "" {
		return ""
	}
	h := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(h[:])
}

func ComparePassword(hashed, plain string) bool {
	return hashed == HashPassword(plain)
}


