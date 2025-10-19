package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	var n int = 10

	wg.Add(n)

	for i := 1; i <= n; i++ {
		var name string = "I am " + fmt.Sprint(i) + " goroutine"
		go CountWithGouroutines(7, name, &wg)
	}

	wg.Wait()
}

func CountWithGouroutines(count int, name string, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < count; i++ {
		fmt.Println(i, name)
	}
}
