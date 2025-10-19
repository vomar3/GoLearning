package main

import "fmt"

// Slice = часть массива, динамически изменяемый размер

func main() {
	var menu = [4]string{"Muffin", "Brownie", "Croissant", "Cookie"}
	fmt.Println(menu)

	// 1 часть включается, 2 не включается
	slice := menu[1:2]
	fmt.Println(slice)

	// All
	slice = menu[:]
	fmt.Println(slice)

	// Все со 2 элемента
	slice = menu[2:]
	fmt.Println(slice)

	// С начало до 3 (не вкл)
	slice = menu[:3]
	fmt.Print(slice)
}
