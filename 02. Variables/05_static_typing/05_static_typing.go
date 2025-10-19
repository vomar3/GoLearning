package main

import "fmt"

func main() {
	// Explicit type declaration
	var cupsQty int = 3

	// Implicit type declaration
	var wasProcessed = true

	//cupsQty = "four" // Очев нельзя
	fmt.Println("Number of cups:", cupsQty)
	fmt.Println(wasProcessed)
}
