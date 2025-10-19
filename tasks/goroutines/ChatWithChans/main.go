package main

import (
	"fmt"
	"time"
)

func main() {
	messages := make(chan string)
	go Alice(messages)

	Bob(messages)
}

func Alice(messages chan<- string) {
	messages <- "Hello, Bob"
	messages <- "How r u?"
	messages <- "Do u want to be a friends?"
	close(messages)
}

func Bob(messages <-chan string) {
	for message := range messages {
		time.Sleep(1 * time.Second)
		fmt.Println(message)
		time.Sleep(1 * time.Second)

		switch message {
		case "Hello, Bob":
			fmt.Println("Hello")
		case "How r u?":
			fmt.Println("I'm ok")
		case "Do u want to be a friends?":
			fmt.Println("Yeah, why not")
		default:
			fmt.Println("Bye")
			return
		}
	}
}
