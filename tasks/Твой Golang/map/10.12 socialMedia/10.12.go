package main

import (
	"fmt"
	"slices"
	"sort"
)

func main() {
	friendsData := map[string][]string{
		"Алексей":  {"Иван", "Сергей", "Елена"},
		"Иван":     {"Алексей", "Дмитрий", "Мария"},
		"Сергей":   {"Алексей", "Елена"},
		"Дмитрий":  {"Иван", "Елена", "Ольга"},
		"Елена":    {"Алексей", "Сергей", "Дмитрий"},
		"Мария":    {"Иван", "Ольга"},
		"Ольга":    {"Дмитрий", "Мария"},
		"Анна":     {"Петр"},
		"Петр":     {"Анна", "Сергей"},
		"Светлана": {"Иван", "Елена"},
	}

	fmt.Print("Количество друзей:\n")

	frCount := countFriends(friendsData)
	users := sortMap(friendsData)

	for _, j := range users {
		fmt.Printf("%s: %d\n", j, frCount[j])
	}

	var name1 string = "Иван"
	var name2 string = "Елена"

	fmt.Printf("Общие друзья между пользователями %s и %s: ", name1, name2)
	var names []string
	names = commonFriends(friendsData, name1, name2)

	for i, value := range names {
		fmt.Printf("%s", value)

		if i == len(names)-1 {
			fmt.Printf(".")
		} else {
			fmt.Printf(", ")
		}
	}

	fmt.Println()

	names = []string{}
	var count int

	names, count = (mostPopularUsers(friendsData))

	for i, value := range names {
		fmt.Printf("%s", value)

		if i == len(names)-1 {
			fmt.Printf(" (количество друзей: %d).", count)
		} else {
			fmt.Printf(", ")
		}
	}
}

func countFriends(friends map[string][]string) map[string]int {
	var count = make(map[string]int, len(friends))

	for key, value := range friends {
		count[key] = len(value)
	}

	return count
}

func commonFriends(friends map[string][]string, firstName string, secondName string) []string {
	var answer []string

	if !(exist(friends, firstName) && exist(friends, secondName)) {
		answer = append(answer, "Никого нет")
		return answer
	}

	checkNames := map[string]int{}

	for _, name := range friends[firstName] {
		checkNames[name]++
	}

	for _, name := range friends[secondName] {
		checkNames[name]++
	}

	for name, value := range checkNames {
		if value == 2 {
			answer = append(answer, name)
		}
	}

	sort.Strings(answer)

	return answer
}

func mostPopularUsers(friends map[string][]string) ([]string, int) {
	var answerSlice []string
	var maxi int = 0

	for key, value := range friends {
		if len(value) > maxi {
			answerSlice = []string{}
			answerSlice = append(answerSlice, key)
			maxi = len(value)
		} else if len(value) == maxi {
			answerSlice = append(answerSlice, key)
		}
	}

	sort.Strings(answerSlice)

	return answerSlice, maxi
}

func exist(friends map[string][]string, name string) bool {
	// checkName == true, if name contains in friends
	if _, checkName := friends[name]; checkName {
		return true
	}

	return false
}

func sortMap(friends map[string][]string) []string {
	result := make([]string, 0, len(friends))

	for key := range friends {
		result = append(result, key)
	}

	slices.Sort(result)

	return result
}
