package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	//"math/rand/v2"
)

func main() {
	//var min, max int = 10, 50

	// Psevdo-random numbers
	//random := rand.IntN(max-min+1) + min
	//fmt.Println(random)

	// Crypto-random numbers
	n, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		log.Fatalf("Ошибка генерации случайного числа: %s", err.Error())
	}

	fmt.Println("Случайное число: ", n.Int64())
}
