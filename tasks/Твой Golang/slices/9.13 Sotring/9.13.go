package main

import (
	"fmt"
	"slices"
)

func main() {
	num := [][]int{
		{3, 1, 4},
		{1, 5, 9},
		{2, 6, 5},
		{0},
	}

	magicSort(num)
}

func magicSort(numbers [][]int) {
	slices.SortFunc(numbers, func(a, b []int) int {
		var sum1, sum2 int

		for _, value := range a {
			sum1 += value
		}

		for _, value := range b {
			sum2 += value
		}

		return sum1 - sum2
	})

	for _, value := range numbers {
		slices.SortFunc(value, func(a, b int) int {
			if (a&1 != 0 && b&1 != 0) || a == 0 {
				return -1
			} else if (a&1 != 0 && b&1 == 0) || b == 0 {
				return 1
			}

			return 0
		})

		slices.SortFunc(value, func(a, b int) int {
			if (a&1 == b&1) && a != 0 && b != 0 {
				return b - a
			}

			return 0
		})
	}

	fmt.Println(numbers)
}
