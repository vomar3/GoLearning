package main

import (
	"errors"
	"fmt"
)

func main() {
	slice := []int{1, 2, 3, 4, 5}
	answer, err := Max(slice)

	if err != nil {
		fmt.Printf("Error: %s", err)
	} else {
		fmt.Println(answer)
	}
}

func Max(numbers []int) (int, error) {
	if len(numbers) == 0 || numbers == nil {
		return 0, errors.New("slice is nil or empty")
	}

	max := numbers[0]

	for _, value := range numbers {
		if max < value {
			max = value
		}
	}

	return max, nil
}
