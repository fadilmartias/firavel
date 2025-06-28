package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateRandomToken(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
