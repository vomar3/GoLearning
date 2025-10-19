package main

import (
	"fmt"
	"math"
)

func main() {
	// Есть 2 способа, со скобками и без скобок, но при сохранении скобки удаляются
	x := 2

	if x < 10 {
		fmt.Println("Hello")
	} else {
		fmt.Println("World")
	}

	if x < 10 {
		fmt.Println("KU")
	}

	// Сначала объявляется переменная, потом ;
	// После ; идет проверка условия if
	if v := math.Pow(2, 3); v < 10 {
		println(52)
	}

	if x < 10 {

	} else if x < 6 {

	} else {

	}
}
