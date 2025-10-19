package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func main() {
	myStr, err := GetInput()

	if err != nil {
		fmt.Println("Error:", err)
		return
	} else {
		a, b, c, d := CountCharacters(myStr)
		DisplayResults(a, b, c, d)
	}
}

func GetInput() (string, error) {
	fmt.Println("Введите строчку для анализа")
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	result := strings.TrimSpace(input)

	if err != nil {
		return "", errors.New("Ошибка ввода строки")
	} else if len(result) == 0 {
		return "", errors.New("Введена нулевая строка, анализ невозможен")
	}

	return result, nil
}

func CountCharacters(text string) (letters, digits, spaces, punctuation int) {
	for _, val1 := range text {
		if unicode.IsPunct(val1) {
			punctuation++
		} else if unicode.IsDigit(val1) {
			digits++
		} else if unicode.IsSpace(val1) {
			spaces++
		} else if unicode.IsLetter(val1) {
			letters++
		}
	}

	return
}

func DisplayResults(letters, digits, spaces, punctuation int) {
	fmt.Printf("Количество букв: %d\n", letters)
	fmt.Printf("Количество цифры: %d\n", digits)
	fmt.Printf("Количество пробелов: %d\n", spaces)
	fmt.Printf("Количество знаков препинания: %d\n", punctuation)
}
