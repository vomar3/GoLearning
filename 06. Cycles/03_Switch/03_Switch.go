package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	// В GO break ставится автоматически в конце case
	// Так же в GO не только целые числа могут быть кейсами
	// Тут можно юзать вообще все

	fmt.Print("Go runs on ")

	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("macOS.")
	case "linux":
		fmt.Println("Linux.")
	default:
		// freebsd, openbsd,
		// plan9, windows...
		fmt.Printf("%s.\n", os)
	}

	// Вот пример с временем
	fmt.Println("When's Saturday?")
	today := time.Now().Weekday()
	fmt.Println(time.Now(), time.Now().Weekday())

	switch time.Saturday {
	case today + 0:
		fmt.Println("Today.")
	case today + 1:
		fmt.Println("Tomorrow.")
	case today + 2:
		fmt.Println("In two days.")
	default:
		fmt.Println("Too far away.")
	}

	// А тут вообще нет условия на свитч, тут просто много кейсов, заменяющих if/else
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("Good morning!")
	case t.Hour() < 17:
		fmt.Println("Good afternoon.")
	default:
		fmt.Println("Good evening.")
	}
}
