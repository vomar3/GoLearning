package main

import "fmt"

func greet(coffeeShopName string) {
	fmt.Println("Welcome to the coffee shop", coffeeShopName)
}

func main() {
	greet("Brew & Beans")
	greet("Coffee & Milk")
}
