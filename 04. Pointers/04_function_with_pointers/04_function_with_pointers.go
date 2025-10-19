package main

import "fmt"

func calculatePrice(price *float64, discountRate float64) {
	*price = *price - (*price * discountRate)
}

func main() {
	var coffeePrice float64 = 5.00
	var discount float64 = 0.10
	fmt.Println("Basic coffee price:", coffeePrice)

	calculatePrice(&coffeePrice, discount)
	fmt.Println("Price with discount", coffeePrice)
}
