package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Сервер запущен")

	if err := http.ListenAndServe(":9091", nil); err != nil {
		fmt.Println("Ошибка, сервер не поднялся")
	}

	fmt.Println("Сервер закрыт")
}
