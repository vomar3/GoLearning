package main

// Если processData выполняется больше 5 секунд, то надо сделать скип (вернуть 0)
// Надо как можно обработалось данных параллельно

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func processData(val int) int {
	time.Sleep(time.Duration(rand.IntN(10)) * time.Second)
	return val * 2
}

func main() {
	in := make(chan int)
	out := make(chan int)

	go func() {
		for i := range 10 {
			in <- i
		}

		close(in)
	}()

	now := time.Now()
	processParallel(in, out, 5)

	for val := range out {
		fmt.Println(val)
	}

	fmt.Println(time.Since(now))
}

func processParallel(in <-chan int, out chan<- int, numWorkers int) {
	wg := &sync.WaitGroup{}

	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for v := range in {
				out <- processData(v)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()
}
