package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
)

func main() {
	for range 100 {
		random = rand.IntN(100) + 1
		guesses = 0
		result := play()
		if result != random {
			fmt.Printf("Неверный ответ. Было загадано число %d, а в ответе получили число %d", random, result)
			os.Exit(-1)
		}
	}
}

var guesses int
var random int

func guess(num int) (int, error) {
	if guesses >= 6 {
		return 0, errors.New("too many attempts")
	}
	guesses++
	if num > random {
		return -1, nil
	}
	if num < random {
		return 1, nil
	}
	return 0, nil
}

func play() int {
	low := 1
	high := 100

	for {
		mid := (low + high) / 2

		number, err := guess(mid)

		if err != nil {
			return mid
		}

		switch number {
		case 0:
			return mid
		case -1:
			high = mid - 1
		case 1:
			low = mid + 1
		}
	}
}
