package main

import (
	"errors"
	"fmt"
)

type Warrior struct {
	Name string
}

type Mage struct {
	Name string
}

type Archer struct {
	Name string
}

func (w Warrior) Attack() string {
	return fmt.Sprintf("Воин %s бьет мечом.", w.Name)
}

func (m Mage) Attack() string {
	return fmt.Sprintf("Маг %s колдует огненный шар.", m.Name)
}

func (a Archer) Attack() string {
	return fmt.Sprintf("Лучник %s выпускает град стрел.", a.Name)
}

func NewWarrior(name string) (*Warrior, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}

	return &Warrior{
		Name: name,
	}, nil
}

func NewMage(name string) (*Mage, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}

	return &Mage{
		Name: name,
	}, nil
}

func NewArcher(name string) (*Archer, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}

	return &Archer{
		Name: name,
	}, nil
}

type Character interface {
	Attack() string
}

func Fight(c []Character) {
	for _, value := range c {
		fmt.Println(value.Attack())
	}
}

func main() {

}
