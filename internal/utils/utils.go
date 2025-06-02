package utils

import (
	"math/rand"
	"strings"
)

func GenerateUniqueCode() string {
	return strings.ToUpper(randomString(6))
}

func randomString(n int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
