package main

import (
	"fmt"
	"maps"
)

func main() {
	new_map := map[string][]int{
		"a": {1, 2, 3},
		"b": {4, 5, 6},
	}

	RemoveSlicesBySum(new_map)
	fmt.Println(new_map)
}

func RemoveSlicesBySum(my_Map map[string][]int) {
	maps.DeleteFunc(my_Map, func(key string, value []int) bool {
		var sum int = 0
		for _, v := range value {
			sum += v
		}

		return sum > 6
	})
}
