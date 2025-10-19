package main

import (
	"fmt"
)

func main() {
	// Константа может подстроиться, но имеет тип int
	const rewardPoints = 10
	var points int = 10

	// int
	fmt.Printf("Default type of rewardPoints is %T\n", rewardPoints)

	var totalRewardPoints float64 = 150.3

	// const могут адаптироваться под другой тип
	totalRewardPoints = totalRewardPoints + rewardPoints

	// Уже надо приводить типы
	totalRewardPoints = totalRewardPoints + float64(points)
}
