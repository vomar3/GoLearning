package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Customer struct {
	ID    int
	Items int
}

func main() {
	var wg sync.WaitGroup
	queue := make(chan Customer)

	wg.Add(3)
	for i := 1; i <= 3; i++ {
		go func(cashierID int) {
			defer wg.Done()
			Cashier(cashierID, queue)
		}(i)
	}

	var customersWg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		customersWg.Add(1)
		go func(customerID int) {
			defer customersWg.Done()
			customer := MakeBuyer(customerID)
			queue <- customer
		}(i)
	}

	go func() {
		customersWg.Wait()
		close(queue)
	}()

	wg.Wait()
	fmt.Println("Магазин закрыт!")
}

func MakeBuyer(i int) Customer {
	randomTime := rand.Float64() * 5.0
	goods := rand.Intn(15) + 1

	time.Sleep(time.Duration(randomTime * float64(time.Second)))

	customer := Customer{ID: i, Items: goods}

	fmt.Printf("Покупатель-%d пришел (%d товаров)\n", customer.ID, customer.Items)

	return customer
}

func Cashier(id int, queue <-chan Customer) {
	for customer := range queue {
		fmt.Printf("Кассир-%d обслуживает Покупателя-%d (%d товаров)\n",
			id, customer.ID, customer.Items)

		serviceTime := time.Duration(customer.Items*300) * time.Millisecond
		time.Sleep(serviceTime)

		fmt.Printf("Кассир-%d закончил с Покупателем-%d\n", id, customer.ID)
	}
}
