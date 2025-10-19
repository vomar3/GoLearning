package main

import (
	"fmt"
	"slices"
)

// Создание динамических массивов

func main() {
	a := make([]int, 5)
	printSlice("a", a)

	b := make([]int, 0, 5)
	printSlice("b", b)

	c := b[:2]
	printSlice("c", c)

	d := c[2:5]
	printSlice("d", d)

	d[0] = 1
	d[1] = 2
	d[2] = 3

	printSlice("d", d)

	d = slices.Delete(d, 1, 3)
	printSlice("d", d)
}

func printSlice(s string, x []int) {
	fmt.Printf("%s len=%d cap=%d %v\n",
		s, len(x), cap(x), x)
}
