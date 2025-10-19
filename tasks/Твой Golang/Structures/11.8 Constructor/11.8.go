package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	MinAge                = 15
	MaxAge                = 80
	MinGrade              = 1
	MaxGrade              = 5
	GreatAge              = 30
	MinGradeAfterGreatAge = 3
)

var (
	ErrEmptyName         = errors.New("name cannot be empty")
	ErrTooYoung          = errors.New("too young")
	ErrTooOld            = errors.New("too old")
	ErrGradeOutOfRange   = errors.New("grade out of range")
	ErrTooLowGradeForAge = errors.New("too low grade for age")
	ErrIncorrectEmail    = errors.New("incorrect email")
)

type Student struct {
	Name  string
	Age   int
	Grade int
	Email string
}

func main() {
	student, err := NewStudent("Dima", 35, 4, "dima@mail.ru")

	if err != nil {
		os.Exit(1)
	} else {
		fmt.Println(*student)
	}
}

func NewStudent(name string, age, grade int, email string) (*Student, error) {
	if name == "" {
		return nil, ErrEmptyName
	}

	if age < MinAge {
		return nil, ErrTooYoung
	} else if age > MaxAge {
		return nil, ErrTooOld
	}

	if grade < MinGrade || grade > MaxGrade {
		return nil, ErrGradeOutOfRange
	}

	if age > GreatAge {
		if grade < MinGradeAfterGreatAge {
			return nil, ErrTooLowGradeForAge
		}
	}

	if !strings.Contains(email, "@") {
		return nil, ErrIncorrectEmail
	}

	return &Student{
		Name:  name,
		Age:   age,
		Grade: grade,
		Email: email,
	}, nil
}
