package main

import "fmt"

func getDrinkInfo(name string, drink string) {
	fmt.Printf("%s's favorite drink is %s", name, drink)
}

func main() {
	getDrinkInfo("Bogdan", "Cappuccino")
}
