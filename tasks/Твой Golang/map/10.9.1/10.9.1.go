package main

import (
	"fmt"
)

func main() {
	new_map := map[string]int{
		"a": 2,
		"b": 4,
	}

	new_map2 := map[string]int{
		"b": 6,
		"a": 3,
	}

	fmt.Println(mergeMaps(new_map, new_map2))
}

func mergeMaps(map1, map2 map[string]int) map[string]int {
	answer := make(map[string]int)

	for key, value := range map1 {
		answer[key] += value
	}

	for key, value := range map2 {
		answer[key] += value
	}

	return answer
}
