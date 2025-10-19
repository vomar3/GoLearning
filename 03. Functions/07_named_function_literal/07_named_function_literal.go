package main

import "fmt"

func main() {
	var taxRate float64 = 0.10

	// named function literal
	// Мы не можем создавать обычные функции в других функциях
	// Но можем создать функцию без имени с параметрами и присваивать переменной
	// Тип = функция
	calculateTax := func(amount float64) float64 {
		return amount * taxRate
	}

	var subtotal float64 = 25.00
	// Вот этот прикольчик работает
	var tax float64 = calculateTax(subtotal)
	total := subtotal + tax

	fmt.Printf("Total amount to pay: $%.2f\n", total)
}
