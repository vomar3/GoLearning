package main

// Можно полностью убрать WG, и будет работать так же (reader удалить go)

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func writer() <-chan int {
	ch := make(chan int)

	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			ch <- i
		}

		close(ch)

	}()

	return ch
}

func doubler(ch <-chan int) <-chan int {
	sch := make(chan int)

	go func() {
		defer wg.Done()
		for val := range ch {
			sch <- (val * 2)
			time.Sleep(500 * time.Millisecond)
		}

		close(sch)
	}()

	return sch
}

func reader(ch <-chan int) {
	go func() {
		defer wg.Done()
		for val := range ch {
			fmt.Println(val)
		}
	}()
}

func main() {
	wg.Add(3)
	reader(doubler(writer()))

	wg.Wait()
}
