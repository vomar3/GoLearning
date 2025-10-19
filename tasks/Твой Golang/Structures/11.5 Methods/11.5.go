package main

import (
	"fmt"
	"math"
)

type Student struct {
	Name   string
	Grades []int
}

func (s Student) AverageGrade() float64 {
	if len(s.Grades) == 0 {
		return 0
	}

	var sum, count int = 0, 0
	for _, val := range s.Grades {
		sum += val
		count++
	}

	return math.Round((float64(sum)/float64(count))*10) / 10
}

func (s Student) Info() string {
	return fmt.Sprintf("Студент %s, средняя оценка: %.1f.", s.Name, s.AverageGrade())
}
