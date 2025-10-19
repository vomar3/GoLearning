package main

import "fmt"

func main() {
	s := [][]int{
		{2, 4, 5, 9},
		{6, 1, 2, 7, 3, 5},
		{4, 7, 4, 8, 0},
	}

	s = replaceEvenOnEvenIndices(s)
	fmt.Println(s)
}

func replaceEvenOnEvenIndices(numbers [][]int) [][]int {
	answer := make([][]int, len(numbers))
	for i, array := range numbers {
		answer[i] = make([]int, len(array))
		copy(answer[i], array)
	}

	for i, array := range numbers {
		for j, value := range array {
			if value%2 == 0 && j%2 == 0 {
				answer[i][j] = 0
			}
		}
	}

	return answer
}
