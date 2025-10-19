package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

var names = []string{
	"Dima",
	"Vanya",
	"Tanya",
	"Max",
	"Anna",
	"Petya",
	"Liza",
	"Andrew",
	"Anton",
	"Sasha",
}

type Order struct {
	ID           int
	CustomerName string
	CookTime     time.Duration
	DeliveryTime time.Duration
}

func main() {
	queue := make(chan Order)
	orders := make(chan Order)

	wgOrders := sync.WaitGroup{}
	for i := 1; i <= 15; i++ {
		random := rand.IntN(10)

		wgOrders.Add(1)
		go func() {
			defer wgOrders.Done()
			order := MakeOrder(i, names[random])
			queue <- order
		}()
	}

	wgCookers := sync.WaitGroup{}
	for i := 1; i <= 4; i++ {
		wgCookers.Add(1)
		go func() {
			defer wgCookers.Done()
			CookPizza(i, queue, orders)
		}()
	}

	wgDelivery := sync.WaitGroup{}
	for i := 1; i <= 2; i++ {
		wgDelivery.Add(1)
		go func() {
			defer wgDelivery.Done()
			DeliveryPizza(i, orders)
		}()
	}

	go func() {
		wgOrders.Wait()
		close(queue)
	}()

	go func() {
		wgCookers.Wait()
		close(orders)
	}()

	wgDelivery.Wait()
	fmt.Println("Пиццерия закрылась")
}

func MakeOrder(id int, name string) Order {
	random := rand.IntN(100) + 1
	time.Sleep(time.Duration(random) * time.Second)

	cook := rand.IntN(3) + 1
	delivery := rand.IntN(4) + 2

	order := Order{
		ID:           id,
		CustomerName: name,
		CookTime:     time.Duration(cook),
		DeliveryTime: time.Duration(delivery),
	}

	fmt.Printf("Заказ-%d от %s (готовить %ds, доставлять %ds)\n", order.ID, order.CustomerName, order.CookTime, order.DeliveryTime)

	return order
}

func CookPizza(id int, queue <-chan Order, orders chan<- Order) {
	for order := range queue {
		fmt.Printf("Повар-%d готовит заказ %d\n", id, order.ID)

		time.Sleep(order.CookTime * time.Second)

		fmt.Printf("Повар-%d закончил заказ %d\n", id, order.ID)

		orders <- order
	}
}

func DeliveryPizza(id int, orders <-chan Order) {
	for delivery := range orders {
		fmt.Printf("Доставщик-%d везет заказ %d\n", id, delivery.ID)

		time.Sleep(delivery.DeliveryTime * time.Second)

		fmt.Printf("Доставщик-%d доставил заказ %d\n", id, delivery.ID)

		time.Sleep(time.Second * 1)
		fmt.Printf("Доставщик-%d едет в пиццерию после доставки заказа %d\n", id, delivery.ID)
	}
}
