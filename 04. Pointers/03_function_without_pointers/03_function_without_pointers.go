package main

import "fmt"

func calculatePrice(price float64, discountRate float64) float64 {
	updatedPrice := price - (price * discountRate)
	return updatedPrice
}

func main() {
	var coffeePrice float64 = 5.00
	var discount float64 = 0.10
	fmt.Println("Basic coffee price:", coffeePrice)

	var newPrice float64 = calculatePrice(coffeePrice, discount)
	fmt.Println("Price with discount", newPrice)
}
