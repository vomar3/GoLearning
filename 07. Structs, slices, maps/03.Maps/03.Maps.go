package main

import "fmt"

type Vertex struct {
	Lat, Long float64
}

var m map[string]Vertex

func main() {
	m = make(map[string]Vertex)
	m["Bell Labs"] = Vertex{
		40.68433, -74.39967,
	}
	fmt.Println(m["Bell Labs"])

	b := make(map[string]int)

	b["Answer"] = 42
	fmt.Println("The value:", b["Answer"])

	b["Answer"] = 48
	fmt.Println("The value:", b["Answer"])

	delete(b, "Answer")
	fmt.Println("The value:", b["Answer"])

	v, ok := b["Answer"]
	fmt.Println("The value:", v, "Present?", ok)
}
