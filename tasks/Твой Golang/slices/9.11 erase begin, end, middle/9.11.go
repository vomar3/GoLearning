package main

import (
	"fmt"
	"slices"
)

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 11}
	fmt.Println(DeletingFromSlice(numbers))
}

func DeletingFromSlice(numbers []int) []int {
	var length int = len(numbers)
	var erasedLast, erasedSecond bool = false, false

	if length == 0 {
		return make([]int, 0)
	}

	if numbers[length-1] > 10 {
		numbers = numbers[:(length - 1)]
		length--
		erasedLast = true
	}

	if length >= 3 && cap(numbers) > 5 {
		numbers = append(numbers[:2], numbers[3:]...)
		length--
		erasedSecond = true
	}

	if erasedLast && erasedSecond && length >= 1 {
		numbers = numbers[1:]
		length--
	}

	numbers = slices.Clip(numbers)

	return numbers
}
