package main

import (
	"errors"
	"fmt"
	"slices"
)

func main() {
	slice1 := []int{1, 2, 3, 3, 3}
	slice2 := []int{3, 4, 1, 52, 54, 3, 100, 12}

	answer, err := intersectSlices(slice1, slice2)
	if err != nil {
		fmt.Printf("Erorr: %s\n", err)
	} else {
		fmt.Println(answer)
	}
}

func intersectSlices(firstSlice, secondSlice []int) ([]int, error) {
	if firstSlice == nil || secondSlice == nil {
		return nil, errors.New("slices cannot be nil")
	}

	answer := []int{}

	if len(firstSlice) > len(secondSlice) {
		for _, value := range secondSlice {
			if slices.Contains(firstSlice, value) {
				answer = append(answer, value)
			}
		}
	} else {
		for _, value := range firstSlice {
			if slices.Contains(secondSlice, value) {
				answer = append(answer, value)
			}
		}
	}

	return answer, nil
}
