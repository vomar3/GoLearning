package main

import "fmt"

func estimateBrewTime(cups int, secondsPerCup int) (totalTimeSeconds int) {
	totalTimeSeconds = cups * secondsPerCup
	return // naked return
}

func main() {
	var time int = estimateBrewTime(12, 20)
	fmt.Println(time)
}
