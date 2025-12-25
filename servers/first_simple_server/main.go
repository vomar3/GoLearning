package main

import (
	"fmt"
	"net/http"
	"time"
)

func payHandler(w http.ResponseWriter, r *http.Request) {
	str := "Новый платеж обработан!"
	b := []byte(str)

	_, err := w.Write(b) // http ответ
	if err != nil {
		fmt.Println("во время записи HTTP ответа произошла ошибка:", err.Error())
	} else {
		fmt.Println("Корректно совершена оплата")
	}
}

func cancelHandler(w http.ResponseWriter, r *http.Request) {
	str := "Оплата отменена!"
	b := []byte(str)

	_, err := w.Write(b) // http ответ
	if err != nil {
		fmt.Println("во время записи HTTP ответа произошла ошибка:", err.Error())
	} else {
		fmt.Println("Корректно отменена оплата")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	str := "Hello, world!"
	b := []byte(str)

	_, err := w.Write(b) // http ответ
	if err != nil {
		fmt.Println("во время записи HTTP ответа произошла ошибка:", err.Error())
	} else {
		fmt.Println("Корретно отработал HTTP запрос")
	}
}

func handlerSleep(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)

	str := "HTTP response!"
	b := []byte(str)

	_, err := w.Write(b) // http ответ
	if err != nil {
		fmt.Println("во время записи HTTP ответа произошла ошибка:", err.Error())
	} else {
		fmt.Println("Корретно отработал HTTP запрос")
	}
}

func main() {
	// Заданы обработчики
	// http.HandleFunc("/", handler) // пустой pattern, пример: youtube.com/
	// Так же называется корневой pattern (endPoint)
	http.HandleFunc("/default", handler)      // pattern default
	http.HandleFunc("/pay", payHandler)       // pattern pay
	http.HandleFunc("/cancel", cancelHandler) // pattern cancel
	http.HandleFunc("/sleep", handlerSleep)

	// Каждый запрос делается в отдельной горутине (под каждый хендлер - отдельная горутина)

	fmt.Println("Запуск сервера")
	err := http.ListenAndServe(":9091", nil) // Запускается сервер и блокирует выполнение программы
	if err != nil {
		fmt.Println("Произошла ошибка:", err.Error())
	}
	fmt.Println("Сервер закрыт")
}
