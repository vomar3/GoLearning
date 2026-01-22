package random

import (
	"math/rand"
)

func Rand(length int) string {
	symbols := "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	answer := make([]byte, length)

	for i := range answer {
		answer[i] = symbols[rand.Intn(len(symbols))]
	}

	return string(answer)
}
