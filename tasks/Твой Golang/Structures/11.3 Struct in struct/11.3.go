package main

import (
	"fmt"
)

// User представляет пользователя
type User struct {
	ID      int
	Name    string
	Email   string
	Phone   string
	Address Address
	Cart    []CartItem
}

// Address представляет адрес пользователя
type Address struct {
	Street     string
	City       string
	PostalCode string
}

// CartItem представляет элемент в корзине
type CartItem struct {
	Product  Product
	Quantity int
}

// Product представляет продукт в корзине
type Product struct {
	ID          int
	Name        string
	Description string
	Price       int
	Category    string
	Brand       string
	Rating      float64
	Reviews     int
}

func printInfo(user User) {
	fmt.Printf("Покупатель %s. Телефон: %s. Адрес: г. %s, %s.\n", user.Name, user.Phone, user.Address.City, user.Address.Street)

	var isBuyer bool = false
	var highPrice = []string{}
	var sum int = 0
	for _, val := range user.Cart {
		if val.Product.Category == "Электроника" {
			isBuyer = true
		}

		if val.Product.Price > 10000 {
			highPrice = append(highPrice, val.Product.Name)
		}

		sum += val.Product.Price * val.Quantity
	}

	fmt.Printf("Пользователь ")
	if isBuyer {
		fmt.Printf("является ")
	} else {
		fmt.Printf("не является ")
	}
	fmt.Printf("покупателем электроники.\n")

	fmt.Printf("Товары в корзине, где цена 10000 и более: ")
	if len(highPrice) == 0 {
		fmt.Printf("отсутствуют.\n")
	} else {
		for i, val := range highPrice {
			if i != len(highPrice)-1 {
				fmt.Printf("%s, ", val)
			} else {
				fmt.Printf("%s.\n", val)
			}
		}
	}

	fmt.Printf("Общая сумма покупки: %d руб.", sum)
}
