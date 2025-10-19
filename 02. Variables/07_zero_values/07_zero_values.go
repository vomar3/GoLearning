package main

import "fmt"

func main() {
	// По умолчанию не мусор, а зануление
	var totalOrders int
	var customerName string
	var isOrderReady bool

	fmt.Println("Total orders:", totalOrders)    // 0
	fmt.Println("Customer name:", customerName)  // ""
	fmt.Println("Order is ready:", isOrderReady) // false
}
