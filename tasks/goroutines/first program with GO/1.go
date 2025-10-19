// Надо сделать так, чтобы 100 горутин отработали < +-5 секунд + посчитать общее время работы функции randomWait (100 - 500 сек)

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var maxWaitSeconds = 5

func randomWait() int {
	workSeconds := rand.Intn(maxWaitSeconds + 1)

	time.Sleep(time.Duration(workSeconds) * time.Second)

	return workSeconds
}

func main() {
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	start := time.Now()

	totalSeconds := 0
	wg.Add(100)
	for range 100 {
		go func() {
			defer wg.Done()
			sec := randomWait()

			mutex.Lock()
			totalSeconds += sec
			mutex.Unlock()
		}()
	}

	wg.Wait()

	mainSeconds := time.Since(start)
	fmt.Println("main: ", mainSeconds)
	fmt.Println("total: ", totalSeconds)
}
