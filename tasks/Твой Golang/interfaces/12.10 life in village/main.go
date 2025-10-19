package main

import (
	"fmt"
	"myvillageproject/village"
)

func main() {
	myVillage := village.Village{}

	// Создаем жителей деревни
	resident1 := &village.Resident{Name: "Алиса", Age: 30, Married: false, Alive: true, Events: []string{}}
	resident2 := &village.Resident{Name: "Борис", Age: 40, Married: true, Alive: true, Events: []string{}}

	// Создаем животных
	animal1 := &village.Animal{Name: "Бобик", Age: 5, Type: "собака", Alive: true, Events: []string{}}
	animal2 := &village.Animal{Name: "Мурка", Age: 3, Type: "кошка", Alive: true, Events: []string{}}
	animal3 := &village.Animal{Name: "Димон", Age: 5, Type: "кошка", Alive: true, Events: []string{}}

	// Добавляем элементы в деревню
	myVillage.AddElement(resident1)
	myVillage.AddElement(resident2)
	myVillage.AddElement(animal1)
	myVillage.AddElement(animal2)

	// Симуляция обновления деревни на несколько лет
	for i := 0; i < 10; i++ {
		fmt.Printf("Сейчас %d год\n", i+1)
		fmt.Print("------------------------------------------------------------\n")
		myVillage.UpdateAll()
		fmt.Println(myVillage.ShowAllInfo())
		fmt.Print("------------------------------------------------------------\n")

		if i == 7 {
			myVillage.AddElement(animal3)
		}
	}

	/*for i := 0; i < 10; i++ {
		fmt.Printf("Сейчас %d год\n", i+1)
		fmt.Print("------------------------------------------------------------\n")

		resident1.Update()
		resident2.Update()
		animal1.Update()
		animal2.Update()

		fmt.Println(resident1.FlushInfo())
		fmt.Println(resident2.FlushInfo())
		fmt.Println(animal1.FlushInfo())
		fmt.Println(animal2.FlushInfo())

		fmt.Print("------------------------------------------------------------\n")
	}*/
}
