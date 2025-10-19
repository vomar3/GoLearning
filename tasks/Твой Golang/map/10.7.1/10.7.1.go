package main

import (
	"fmt"
	"sort"
)

func main() {
	new_map := map[string]int{
		"Mitchel Resnick":   5,
		"Linus Torvalds":    5,
		"Donald Knuth":      3,
		"Tim Berners-Lee":   5,
		"Bjarne Stroustrup": 5,
	}

	fmt.Println(countVotes(new_map))
}

func countVotes(votes map[string]int) string {
	if len(votes) == 0 {
		return "Кандидаты потерялись."
	}

	var answer string
	var maxi int = 0
	var names = make([]string, 0, len(votes))

	for i := range votes {
		maxi = max(maxi, votes[i])
		names = append(names, i)
	}

	if maxi == 0 {
		return "Все голоса похищены!"
	}

	sort.Strings(names)

	for _, name := range names {
		if votes[name] == maxi {
			if len(answer) != 0 {
				answer += ", "
			}

			answer += name
		}
	}

	return answer
}
