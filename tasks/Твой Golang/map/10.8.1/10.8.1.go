package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	new_map := map[string][]int{
		"a": {1, 2, 3},
		"b": {4, 5, 6},
	}

	new_map2 := map[string][]int{
		"b": {6, 5, 4},
		"a": {3, 1, 2},
	}

	fmt.Println(CompareMaxValues(new_map, new_map2))
}

func CompareMaxValues(map1, map2 map[string][]int) bool {
	result := maps.EqualFunc(map1, map2, func(v1, v2 []int) bool {
		if len(v1) == 0 && len(v2) == 0 {
			return true
		} else if len(v1) == 0 || len(v2) == 0 {
			return false
		}

		return slices.Max(v1) == slices.Max(v2)
	})

	return result
}
