package link

import (
	"math/rand/v2"
)

func GenerateShort(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	link := make([]byte, length)
	for i := range link {
		link[i] = charset[rand.IntN(len(charset))]
	}

	return string(link)
}
