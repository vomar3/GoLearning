package main

import "fmt"

func main() {
	var coffeePrice = 4.50
	fmt.Println("Coffee price:", coffeePrice)

	fmt.Println("Memory adress:", &coffeePrice)

	var pointer *float64 = &coffeePrice
	*pointer = 5.50

	fmt.Println("Updated coffeePrice value:", coffeePrice)
}
