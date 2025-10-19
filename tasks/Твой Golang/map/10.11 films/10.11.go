package main

import (
	"fmt"
	"maps"
	"math"
	"sort"
)

type Pair struct {
	Key   string
	Value float64
}

func main() {
	m := map[string]map[string]float64{
		"Экшен": {
			"Фильм1": 8.52,
			"Фильм2": 6.0,
		},
		"Драма": {
			"Фильм3": 7.524,
			"Фильм4": 7.527,
			"Фильм5": 5.54,
		},
		"Окак": {
			"Фильм7": 5.52,
		},
	}

	printRecommendations(m)
}

func printRecommendations(movies map[string]map[string]float64) {
	var names []string
	var eps float64 = 0.0001

	for key, value := range movies {
		maps.DeleteFunc(value, func(s_key string, val float64) bool {
			if math.Abs(val-eps) < 7.0 {
				return true
			}
			return false
		})

		if len(value) > 0 {
			names = append(names, key)
		}
	}

	sort.Strings(names)

	for _, value := range names {
		fmt.Printf("%s: ", value)

		var KeysAndValues []Pair

		for key, val := range movies[value] {
			KeysAndValues = append(KeysAndValues, Pair{key, val})
		}

		sort.Slice(KeysAndValues, func(i, j int) bool {
			if KeysAndValues[i].Value == KeysAndValues[j].Value {
				return KeysAndValues[i].Key < KeysAndValues[j].Key
			}

			return KeysAndValues[i].Value > KeysAndValues[j].Value
		})

		for key, _ := range KeysAndValues {
			fmt.Printf("%s (%.1f)", KeysAndValues[key].Key, KeysAndValues[key].Value)

			if key != len(KeysAndValues)-1 {
				fmt.Print(", ")
			} else {
				fmt.Print(".")
			}
		}

		fmt.Println()
	}
}
