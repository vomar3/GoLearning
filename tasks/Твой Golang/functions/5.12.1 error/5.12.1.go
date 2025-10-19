package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

func main() {
	value, err := UserProfileToString("Dima", -15)

	if err != nil {
		log.Fatalf("Error: %s", err)
	} else {
		fmt.Println(value)
	}
}

func UserProfileToString(name string, age int) (string, error) {
	if age < 0 {
		return "", errors.New("negative age")
	}

	if name == "" {
		return "", errors.New("empty name")
	}

	result := strings.TrimSpace(name)

	if result == "" {
		return "", errors.New("name cannot contain only spaces")
	}

	answer := fmt.Sprintf("Имя человека: %s, возраст: %d.", result, age)
	return answer, nil
}
