package main

import (
	"fmt"
	"math/rand"
	"time"
)

func randomTimeWork() {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Second)
}

func predictableTimeWork() {
	ch := make(chan struct{})

	go func() {
		randomTimeWork()
		close(ch)
	}()

	select {
	case <-ch:
		fmt.Println("all good")
		return
	case <-time.After(3 * time.Second):
		fmt.Println("error")
		return
	}
}

func main() {
	predictableTimeWork()
}
