package main

import "fmt"

func main() {
	var totalPoint int = 120
	var points int = calculate(25.5)
	fmt.Println(points)

	totalPoint = plus(totalPoint, points)
	fmt.Println("Updated lotalty points:", totalPoint)
}

func plus(calc int, total int) int {
	return calc + total
}

func calculate(spent float64) int {
	var points int = int(spent * 2)
	return points
}
