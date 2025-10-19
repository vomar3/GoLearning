package main

import (
	"fmt"
	"sort"
)

func main() {
	m := map[string]int{
		"banana":     2,
		"apple":      1,
		"grapefruit": 3,
		"cherry":     1,
	}

	n := invertMap(m)
	printMap(n)
}

func invertMap(my_Map map[string]int) map[int][]string {
	var answer = make(map[int][]string)

	for key, value := range my_Map {
		answer[value] = append(answer[value], key)
	}

	return answer
}

func printMap(my_Map map[int][]string) {
	fmt.Println("{")

	var keys []int

	for key := range my_Map {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	for _, key := range keys {
		value := my_Map[key]
		sort.Strings(value)

		fmt.Printf("  %d:", key)
		fmt.Print(" [")

		for i, str := range value {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("\"%s\"", str)
		}
		fmt.Println("],")
	}

	fmt.Println("}")
}
