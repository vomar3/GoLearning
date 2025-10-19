package main

import "fmt"

func main() {
	// Такое использование переменных только внутри функций
	name := "Americano"
	price := 2.99
	ready := true
	count := 5
	var stockCount int64 = 5000

	fmt.Printf("Type of name is: %T\n", name)
	fmt.Printf("Type of price is: %T\n", price)
	fmt.Printf("Type of ready is: %T\n", ready)
	fmt.Printf("Type of count is: %T\n", count)
	fmt.Printf("Type of stockCount is: %T\n", stockCount)
}
