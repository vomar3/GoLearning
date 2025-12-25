package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Payment struct {
	// Описание покупки
	Description string `json:"description"`
	// Сумма покупки
	USD int `json:"usd"`
	// ФИО человека, совершающего покупку
	FullName string `json:"fullName"`
	// Место прописки
	Address string `json:"address"`
	// Время проведения оплаты
	Time time.Time
}

func (p Payment) Println() {
	fmt.Println("Description:", p.Description)
	fmt.Println("USD:", p.USD)
	fmt.Println("FullName:", p.FullName)
	fmt.Println("Address:", p.Address)
}

// Проверка валидации данных из JSON
func (p Payment) Validate() bool {
	if p.USD == 0 {
		return false
	}

	if p.Address == "" {
		return false
	}

	if p.FullName == "" {
		return false
	}

	// и т.д. проверки

	return true
}

var mtx = sync.Mutex{}
var money = 1000
var paymentHistory = make([]Payment, 0)

type HttpResponse struct {
	Money          int       `json:"money"`
	PaymentHistory []Payment `json:"phistory"`
}

func payHandler(w http.ResponseWriter, r *http.Request) {
	// Парсится json
	// Потеря данных != ошибка, если нет данных, то прилетает просто нулевое значение
	var payment Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		fmt.Println("err:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	payment.Time = time.Now()

	payment.Println()

	mtx.Lock()
	if money-payment.USD >= 0 {
		money -= payment.USD
	}

	paymentHistory = append(paymentHistory, payment)

	httpResponse := HttpResponse{
		Money:          money,
		PaymentHistory: paymentHistory,
	}

	// b, err := json.MarshalIndent(httpResponse, "", "	") // Форматированный JSON
	b, err := json.Marshal(httpResponse) // Простой JSON
	if err != nil {
		fmt.Println("err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("err:", err)
		return
	}

	mtx.Unlock()
}

func main() {
	fmt.Println("Сервер запущен")

	http.HandleFunc("/pay", payHandler)

	if err := http.ListenAndServe(":9091", nil); err != nil {
		fmt.Println("Ошибка во время работы HTTP сервера", err)
	}

	fmt.Println("Сервер закрыт")
}
