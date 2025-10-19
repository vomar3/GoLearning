package village

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"
)

var ResidentEvents = []string{
	"Поругался с соседом",
	"Поел черную икру на новый год",
	"Купил машину",
	"Сломал руку",
	"Сломал ногу",
	"Получил повышение на работе",
	"Взял кредит",
	"Купил дом",
	"Сильно заболел",
	"Бахнул пивка",
}

type Resident struct {
	Name    string
	Age     int
	Married bool
	Alive   bool
	Events  []string
}

func NewResident(name string, age int, married, alive bool) (*Resident, error) {
	if name == "" {
		return nil, errors.New("у человека должно быть имя")
	}

	if age < 0 {
		return nil, errors.New("возраст должен быть >= 0")
	}

	return &Resident{
		Name:    name,
		Age:     age,
		Married: married,
		Alive:   alive,
		Events:  make([]string, 0, 1),
	}, nil
}

func (R *Resident) UpdateYear() {
	R.Age++
}

func (R *Resident) Marriage() {
	if R.Married {
		R.Married = false
		R.Events = append(R.Events, "Развод, больше я не в браке.")
	} else {
		R.Married = true
		R.Events = append(R.Events, "Наконец-то, найден спутник в жизни!!!")
	}
}

func (R *Resident) UpdateAlive() {
	R.Alive = false

	R.Events = append(R.Events, fmt.Sprintf("К великому сожалению, наш любимый человек - %s умер на %d году жизни. Почтим память", R.Name, R.Age))
}

func (R *Resident) Update() { // Дописать возможность смерти
	if R.Alive {
		R.Events = nil
		R.UpdateYear()

		countEvents := rand.IntN(3)

		for i := 0; i < countEvents; i++ {
			randomNumber := rand.IntN(10)
			R.Events = append(R.Events, ResidentEvents[randomNumber])
		}

		randomNumber := rand.IntN(10)
		if randomNumber == 5 {
			R.Marriage()
		}

		randomNumber = rand.IntN(50)
		if randomNumber == 27 {
			R.UpdateAlive()
		}

		if R.Events == nil {
			R.Events = append(R.Events, "Скучная жизнь в этом году была ;(")
		}
	} else {

	}
}

func (R *Resident) FlushInfo() string {
	eventsText := "События: нет"

	status := ""
	if R.Married {
		status = "в браке"
	} else {
		status = "холост"
	}

	if len(R.Events) > 0 {
		eventsText = "События:\n" + strings.Join(R.Events, "\n")
	}

	R.Events = nil

	return fmt.Sprintf("Житель %s (возраст: %d), статус: %s\n%s", R.Name, R.Age, status, eventsText)
}

func (R Resident) CheckAlive() bool {
	return R.Alive
}
