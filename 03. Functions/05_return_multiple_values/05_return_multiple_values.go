package main

import "fmt"

func processPay(orderTotal float64, tip float64, paid float64) (float64, float64) {
	total := orderTotal + tip
	change := paid - total

	return total, change
}

func main() {
	var first, second float64 = processPay(6.50, 2.00, 10.00)
	fmt.Println(first, second)
}
