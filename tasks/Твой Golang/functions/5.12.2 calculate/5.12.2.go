package main

import (
	"errors"
	"fmt"
	"log"
	"math"
)

func main() {
	value, err := calculate(10., 0., "divide")

	if err != nil {
		log.Fatalf("Error check: %s", err)
	}

	fmt.Println(value)
}

func calculate(val1, val2 float64, operation string) (float64, error) {
	switch operation {
	case "add":
		return val1 + val2, nil
	case "subtract":
		return val1 - val2, nil
	case "multiply":
		return val1 * val2, nil
	case "divide":
		eps := 0.001
		if math.Abs(val2) < eps {
			return 0, errors.New("division by zero")
		}

		return val1 / val2, nil
	default:
		return 0, errors.New("unknown operation")
	}
}
