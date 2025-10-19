package main

import "fmt"

func main() {
	fmt.Println(filterEven(2, 3, 4, 5, 6, 9, 0, 1, 3, 52))
}

func filterEven(numbers ...int) []int {
	var answer []int

	for _, value := range numbers {
		if value&1 == 0 {
			answer = append(answer, value)
		}
	}

	return answer
}
