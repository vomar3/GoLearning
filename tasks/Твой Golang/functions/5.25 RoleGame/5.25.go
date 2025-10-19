package main

import (
	"fmt"
	"math/rand/v2"
)

func main() {
	attack := DamageBoostDecorator(CriticalHitDecorator(SlowEffectDecorator(Attack)))
	fmt.Println(attack())
}

func Attack() string {
	return "Атака выполнена!"
}

func DamageBoostDecorator(attackFunc func() string) func() string {
	return func() string {
		return "Вам улыбнулась удача, нанесение урона увеличено на 10%!\n" + attackFunc()
	}
}

func CriticalHitDecorator(attackFunc func() string) func() string {
	return func() string {
		if rand.IntN(100) < 25 {
			return "Критический удар! Урон удвоен!\n"
		} else {
			return attackFunc()
		}
	}
}

func SlowEffectDecorator(attackFunc func() string) func() string {
	return func() string {
		return attackFunc() + "Цель замедлена на 2 хода!"
	}
}
