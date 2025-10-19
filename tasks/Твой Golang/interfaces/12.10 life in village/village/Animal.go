package village

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"
)

var AnimalEvents = []string{
	"Поругался с хозяином",
	"Поел черную икру на новый год",
	"Покусал прохожего",
	"Сломал лапу",
	"Подрался с зеркалом",
	"Шлепнулся с крыльца",
	"Получил похвалу от соседки",
	"Отвоевал огород",
	"Сходил к ветеринару",
	"Бахнул пивка",
}

type Animal struct {
	Name   string
	Age    int
	Type   string
	Alive  bool
	Events []string
}

func NewAnimal(name string, age int, Atype string, alive bool) (*Animal, error) {
	if name == "" {
		return nil, errors.New("у животного должно быть имя")
	}

	if age < 0 {
		return nil, errors.New("возраст должен быть >= 0")
	}

	if Atype == "" {
		return nil, errors.New("тип животного должен быть введен")
	}

	return &Animal{
		Name:   name,
		Age:    age,
		Type:   Atype,
		Alive:  alive,
		Events: make([]string, 0, 1),
	}, nil
}

func (A *Animal) UpdateYear() {
	A.Age++
}

func (A *Animal) UpdateAlive() {
	A.Alive = false

	A.Events = append(A.Events, fmt.Sprintf("К великому сожалению, наш любимый питомец  - %s умер на %d году жизни. Почтим память", A.Name, A.Age))
}

func (A *Animal) Update() { // Дописать возможность смерти
	if A.Alive {
		A.Events = nil
		A.UpdateYear()

		countEvents := rand.IntN(3)

		for i := 0; i < countEvents; i++ {
			randomNumber := rand.IntN(10)
			A.Events = append(A.Events, AnimalEvents[randomNumber])
		}

		randomNumber := rand.IntN(10)
		if randomNumber == 5 {
			A.UpdateAlive()
		}

		if A.Events == nil {
			A.Events = append(A.Events, "Скучная жизнь в этом году была ;(")
		}
	}
}

func (A *Animal) FlushInfo() string {
	eventsText := "События: нет"

	if len(A.Events) > 0 {
		eventsText = "События:\n" + strings.Join(A.Events, "\n")
	}

	A.Events = nil

	return fmt.Sprintf("Животное %s (%s, возраст: %d).\n%s", A.Name, A.Type, A.Age, eventsText)
}

func (A Animal) CheckAlive() bool {
	return A.Alive
}
