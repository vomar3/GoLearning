package main

import (
	"errors"
	"fmt"
)

type TagManager struct {
	Tags map[string]struct{}
}

func main() {
	tm := NewTagManager()

	// Добавление тегов
	if err := tm.AddTag("golang"); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	if err := tm.AddTag("programming"); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	if err := tm.AddTag("golang"); err != nil {
		fmt.Printf("Error: %v\n", err) // Ошибка, тег уже существует
	}

	// Проверка существования тегов
	fmt.Println("Тег 'golang' существует:", tm.TagExists("golang")) // true
	fmt.Println("Тег 'python' существует:", tm.TagExists("python")) // false

	// Список тегов
	fmt.Println("Current tags:", tm.ListTags()) // [golang programming]

	// Удаление тегов
	if err := tm.RemoveTag("golang"); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	if err := tm.RemoveTag("golang"); err != nil {
		fmt.Printf("Error: %v\n", err) // Ошибка, тег не существует
	}

	// Список тегов после удаления
	fmt.Println("Current tags after removal:", tm.ListTags()) // [programming]
}

func NewTagManager() *TagManager {
	return &TagManager{
		Tags: make(map[string]struct{}, 5),
	}
}

func (Manager *TagManager) AddTag(tag string) error {
	if _, err := Manager.Tags[tag]; err {
		return errors.New("Tag is alreary exist")
	}

	Manager.Tags[tag] = struct{}{}

	return nil
}

func (Manager *TagManager) RemoveTag(tag string) error {
	if _, err := Manager.Tags[tag]; !err {
		return errors.New("Tag isn't exist")
	}

	delete(Manager.Tags, tag)
	return nil
}

func (Manager *TagManager) TagExists(tag string) bool {
	if _, err := Manager.Tags[tag]; err {
		return true
	}

	return false
}

func (Manager *TagManager) ListTags() []string {
	var answer []string

	for key := range Manager.Tags {
		answer = append(answer, key)
	}

	return answer
}
