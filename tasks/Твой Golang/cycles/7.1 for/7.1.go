package main

import (
	"fmt"
	"math/rand/v2"
)

func main() {
	dice := 9
	rollDice(dice)
}

func rollDice(dice int) {
	roll, rollFirst, rollSecond, count := 0, 0, 0, 0

	for roll != dice {
		rollFirst = rand.IntN(6) + 1
		rollSecond = rand.IntN(6) + 1
		count++

		roll = rollFirst + rollSecond

		if roll != dice {
			fmt.Printf("Выпало %d и %d, в сумме %d, бросаем еще раз.\n", rollFirst, rollSecond, roll)
		}
	}

	fmt.Printf("Выпало %d и %d, в сумме %d, на это ", rollFirst, rollSecond, roll)

	switch {
	case (count%10 == 1 && count != 11):
		fmt.Printf("потребовался %d бросок.\n", count)
	case (count%10 >= 2 && count%10 <= 4 && count != 12 && count != 13 && count != 14):
		fmt.Printf("потребовалось %d броска.\n", count)
	default:
		fmt.Printf("потребовалось %d бросков.\n", count)
	}

}
