package main

import (
	"fmt"
	"strconv"
)

func main() {
	var array = [5]int{3, 8, 1, 8, 1}
	answer := SecretCode(array)
	fmt.Println(answer)
}

func SecretCode(array [5]int) string {
	answer := ""

	var min, max int = 10, 0

	for i := 0; i < 5; i++ {
		if max < array[i] {
			max = array[i]
		}

		if min > array[i] {
			min = array[i]
		}
	}

	answer += strconv.Itoa(min)

	for i := 0; i < 5; i++ {
		if array[i]&1 == 0 {
			answer += string('E')
		} else {
			answer += string('O')
		}

		answer += strconv.Itoa(array[i])
	}

	answer += strconv.Itoa(max)

	return answer
}
