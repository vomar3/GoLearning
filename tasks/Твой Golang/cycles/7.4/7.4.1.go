package main

import (
	"fmt"
)

func main() {
	printTable(3)
}

func printTable(num int) {
	for i := 1; i <= num; i++ {
		for j := 1; j <= num; j++ {
			fmt.Printf("%d x %d = %d\t", i, j, i*j)
		}
		fmt.Println()
	}
}
