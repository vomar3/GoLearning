package main

import (
	"fmt"
)

func main() {
	add := adder(10)     // Создаем замыкание с начальным значением 10
	fmt.Println(add(5))  // Ожидаемый результат: 15
	fmt.Println(add(10)) // Ожидаемый результат: 25
}

func adder(n int) func(x int) int {
	value := n

	return func(x int) int {
		value = value + x
		return value
	}
}
