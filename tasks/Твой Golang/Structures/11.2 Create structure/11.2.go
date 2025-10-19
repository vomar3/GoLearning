package main

import "fmt"

type Event struct {
	Title    string
	Date     string
	Location string
}

func main() {
	user := createGoEvent()
	fmt.Println(user)
}

func createGoEvent() Event {
	var create = Event{
		Title:    "День рождения Golang",
		Date:     "10 ноября 2009",
		Location: "GoogleLand",
	}

	return create
}
