package main

import (
	"fmt"
	"slices"
)

func main() {
	numbers := []int{1, 2, 3, 4, 5}
	fmt.Println(PlayWithSlice(numbers))
}

func PlayWithSlice(numbers []int) []int {
	newNumbers := slices.Clone(numbers)

	var idTenPlus, sum, countEven, countOdd int = -1, 0, 0, 0

	for i, value := range numbers {
		sum += value

		if value > 10 {
			idTenPlus = i + 1
		}

		if value&1 == 0 {
			countEven++
		} else {
			countOdd++
		}
	}

	if idTenPlus != -1 {
		newNumbers = append(newNumbers[:idTenPlus], append([]int{100}, newNumbers[idTenPlus:]...)...)

		fmt.Println(newNumbers)

		sum += 100
	}

	if sum > 100 {
		newNumbers = append(newNumbers, 500)
	}

	if countEven > countOdd {
		newNumbers = append([]int{1000}, newNumbers...)
	}

	return newNumbers
}
