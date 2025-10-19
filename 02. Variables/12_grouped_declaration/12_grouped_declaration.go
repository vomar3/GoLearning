package main

import "fmt"

func main() {
	var coffeeType string = "Latte"
	var quantity int = 3
	var unitPrice float64 = 4.25

	var (
		customerName string = "Bogdan"
		tableNumber  int    = 8
		isReadyToPay bool   = false
	)

	const (
		SizeSmall  = "S"
		SizeMedium = "M"
		SizeLarge  = "L"
	)

	fmt.Println(coffeeType, quantity, unitPrice,
		customerName, tableNumber, isReadyToPay)

	fmt.Printf("%t", isReadyToPay)

}
