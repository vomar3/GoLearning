package main

import (
	"fmt"
)

func main() {
	str := "Зашифруй меня!"
	encodedStr := CaesarCode(str, 1, true)
	fmt.Println(encodedStr)

	decodedStr := CaesarCode(encodedStr, 1, false)
	fmt.Println(decodedStr)
}

func CaesarCode(text string, shift int32, encode bool) string {
	runes := []rune(text)

	if encode {
		for i := range runes {
			runes[i] += shift
		}
	} else {
		for i := range runes {
			runes[i] -= shift
		}
	}

	return string(runes)
}
