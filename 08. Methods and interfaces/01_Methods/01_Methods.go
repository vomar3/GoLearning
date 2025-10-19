package main

import (
	"fmt"
	"math"
)

type Vertex struct {
	X, Y float64
}

// В Go нет классов, поэтому просто создается метод-функция, где в аргументах есть v Vertex
// В этом примере метод Abs имеет приемник типа Vertex с именем v.
func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func main() {
	v := Vertex{3, 4}
	// Использование как своего метода, а не как вызов функции
	fmt.Println(v.Abs())

	f := MyFloat(-math.Sqrt2)
	fmt.Println(f.Abs())
}

type MyFloat float64

func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}
