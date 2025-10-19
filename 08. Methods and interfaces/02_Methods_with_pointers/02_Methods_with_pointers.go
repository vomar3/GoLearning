// pointer receiver

package main

import (
	"fmt"
	"math"
)

type Vertex struct {
	X, Y float64
}

// Просто метод, который можно вызывать, и он НЕ МЕНЯЕТ значение переменной

func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Вызов метода МЕНЯЕТ значение переменной

func (v *Vertex) Scale(f float64) {
	// (*v).X = (*v).X * f  // Так тоже можно, но не нужно
	v.X = v.X * f
	v.Y = v.Y * f
}

func main() {
	v := Vertex{3, 4}
	v.Scale(10) // Go автоматически делает (&v).Scale(10)
	fmt.Println(v.Abs())
}
