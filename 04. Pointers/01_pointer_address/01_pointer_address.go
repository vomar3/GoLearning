package main

import "fmt"

func main() {
	coffee := "Espresso"
	var pointer *string = &coffee // pointer := &coffee

	fmt.Println("Coffee name:", coffee)
	fmt.Println("Memory location:", pointer)
	fmt.Printf("Pointer address: %p\n", pointer)

	fmt.Println("---------------")

	coffeeCopy := coffee

	fmt.Println("Coffee name:", coffeeCopy)
	fmt.Println("Memory location:", &coffeeCopy)
	fmt.Printf("Pointer address: %p\n", &coffeeCopy)
}
