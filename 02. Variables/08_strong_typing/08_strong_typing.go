package main

import "fmt"

func main() {
	// Price of one cup
	price := 4.50 // float64

	// Cups
	quantity := 15 // int

	// total income
	total := price * float64(quantity) // Ну тут тоже очев

	fmt.Printf("Total income: %.2f\n", total)
}
