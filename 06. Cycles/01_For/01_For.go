package main

import "fmt"

func main() {
	var sum int = 0

	for i := 0; i < 10; i++ {
		sum += i
	}

	fmt.Println(sum)

	for sum < 100 {
		sum += sum
	}

	fmt.Println(sum)

	i := 0

	j := 5
	// Нельзя писать так
	// j = i++
	//j = ++i

	j += i + 1
}
