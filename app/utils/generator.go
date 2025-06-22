package utils

import (
	"math/rand"
	"time"
)

// Daftar karakter untuk ID pendek
const shortIDChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateShortID(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = shortIDChars[seededRand.Intn(len(shortIDChars))]
	}
	return string(b)
}
