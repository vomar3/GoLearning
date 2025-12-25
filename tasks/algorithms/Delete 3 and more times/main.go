package main

import (
	"fmt"
)

func main() {
	var array = []int{1, 2, 1, 4, 5, 9, 19, 20, 1}

	array = deleteRepetitions(array)
	fmt.Println(array)
}

func deleteRepetitions(array []int) []int {
	repetMap := make(map[int]int, len(array))
	answer := make([]int, 0, len(array))

	for i := 0; i < len(array); i++ {
		repetMap[array[i]]++
		if repetMap[array[i]] <= 2 {
			answer = append(answer, array[i])
		}
	}

	return answer
}
