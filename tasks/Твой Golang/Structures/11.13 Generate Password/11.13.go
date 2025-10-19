package main

import (
	"crypto/rand"
	"errors"
	"math/big"
)

const (
	MinPasswordLength = 4
	MinPasswordsCount = 1
	MaxPasswordsCount = 50
)

var (
	ErrPasswordLengthTooLow = errors.New("password length too low")
	ErrTooLowPasswordsCount = errors.New("too low passwords count")
	ErrTooBigPasswordsCount = errors.New("too big passwords count")
)

var (
	upperChars   = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	lowerChars   = []rune("abcdefghijklmnopqrstuvwxyz")
	digitChars   = []rune("0123456789")
	specialChars = []rune("!@#$%^&*")
)

func main() {

}

func generatePassword(length int, count int) ([]string, error) {
	if length < MinPasswordLength {
		return nil, ErrPasswordLengthTooLow
	}

	if count < MinPasswordsCount {
		return nil, ErrTooLowPasswordsCount
	} else if count > MaxPasswordsCount {
		return nil, ErrTooBigPasswordsCount
	}

	answer := make([]string, 0, count)
	allChars := append(append(append(upperChars, lowerChars...), digitChars...), specialChars...)

	for j := 0; j < count; j++ {
		password := make([]rune, length)

		password[0] = getRandomRune(upperChars)
		password[1] = getRandomRune(lowerChars)
		password[2] = getRandomRune(digitChars)
		password[3] = getRandomRune(specialChars)

		for i := 4; i < length; i++ {
			password[i] = getRandomRune(allChars)
		}

		shuffleRunes(password)
		answer = append(answer, string(password))
	}

	return answer, nil
}

func getRandomRune(chars []rune) rune {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
	return chars[num.Int64()]
}

func shuffleRunes(slice []rune) {
	for i := len(slice) - 1; i > 0; i-- {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		slice[i], slice[j.Int64()] = slice[j.Int64()], slice[i]
	}
}
