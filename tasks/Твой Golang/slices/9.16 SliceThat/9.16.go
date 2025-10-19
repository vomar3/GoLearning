package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"slices"
)

func main() {
	slice, err := CreateSlice(10)

	if err != nil {
		fmt.Printf("Error %s\n", err)
		os.Exit(1)
	}

	fmt.Println(slice)
	slice = SortByParity(slice)
	fmt.Println(slice)
	slice = MaxSumWithNegative(slice, 3)
	fmt.Println(slice)
	slice = FilterSlice(slice)
	fmt.Println(slice)
}

func CreateSlice(n int) ([]int, error) {
	if n < 0 {
		return nil, errors.New("Попытка создать массив с отрицательной длиной, аварийное завершение программы")
	}

	answer := make([]int, 0, n)

	for range n {
		answer = append(answer, rand.Intn(21)-10)
	}

	return answer, nil
}

func FilterSlice(numbers []int) []int {
	slice := []int{}

	for i := 1; i < len(numbers); i++ {
		if numbers[i] < numbers[i-1] {
			if numbers[i]%2 == 0 || numbers[i]%5 == 0 ||
				numbers[i]%6 == 0 || numbers[i]%9 == 0 {
				slice = append(slice, numbers[i])
			}
		}
	}

	return slice
}

func MaxSumWithNegative(numbers []int, k int) []int {
	if k > len(numbers) {
		return nil
	}

	slice := make([]int, k)
	max := -30

	for i := 0; i <= len(numbers)-k; i++ {
		sum := SumSlice(numbers[i : k+i])
		if sum > max && CheckNegativeNumbers(numbers[i:k+i]) {
			max = sum
			slice = numbers[i : k+i]
		}
	}

	return slice
}

func SumSlice(slice []int) int {
	sum := 0

	for _, val := range slice {
		sum += val
	}

	return sum
}

func CheckNegativeNumbers(slice []int) bool {
	for _, val := range slice {
		if val < 0 {
			return true
		}
	}

	return false
}

func SortByParity(numbers []int) []int {
	fmt.Println("before:", numbers)

	slices.SortFunc(numbers, func(a, b int) int {
		if a%2 == 0 && b%2 != 0 {
			return -1
		} else if a%2 != 0 && b%2 == 0 {
			return 1
		}

		return 0
	})

	slices.SortFunc(numbers, func(a, b int) int {
		if a%2 == 0 && b%2 == 0 {
			return b - a
		} else if a%2 != 0 && b%2 != 0 {
			return a - b
		}

		return 0
	})

	fmt.Println("after:", numbers)

	return numbers
}
