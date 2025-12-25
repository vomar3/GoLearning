package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

var mtx sync.Mutex
var money = 1000 // usd
var bank = 0     // usd

func payHandler(w http.ResponseWriter, r *http.Request) {
	// В хедере зашита служебная инфа, например
	// формат отправленных данных, юзер-агент (кто клиент), токены и т.д.
	// По факту, нужно, чтобы люди не могли получить доступ к чужим данным, например
	/*for k, v := range r.Header {
		fmt.Println("k:", k, "-- v:", v)
	}*/

	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed) // Patch
		return
	}

	fmt.Println("HTTP method:", r.Method) // get, post, delete ...

	httpRequestBody, err := io.ReadAll(r.Body) // byte[] возвращает
	if err != nil {
		// http status codes (почитать)
		// По умолчанию кидается 200
		w.WriteHeader(http.StatusInternalServerError) // 500

		msg := fmt.Sprintf("fail to read HTTP body: %v", err)
		fmt.Println(msg)
		_, err := w.Write([]byte(msg)) // закидываем ответ на сервер
		if err != nil {
			fmt.Println("Fail to write HTTP response", err)
		}
		return
	}

	httpRequestBodyString := string(httpRequestBody)

	paymentAmount, err := strconv.Atoi(httpRequestBodyString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400

		msg := "Fail to convert hhtpBody to integer" + err.Error()
		// msg := fmt.Sprintf("fail to read HTTP body: %v", err) ТАК ЛУЧШЕЫ
		fmt.Println(msg)
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("Fail to write HTTP response", err)
		}
		return
	}

	// Нужен мьютекс, т.к. используются глоабальные переменные, а хендлеры - горутины
	// Есть вероятность, что 2 хендлера на списание денег и откладывание в копилку
	// Вызовутся одновременно => баланс может уйти в минус => блокируется мьютекс
	mtx.Lock()
	if money-paymentAmount >= 0 {
		money -= paymentAmount

		msg := "Оплата прошла успешно, баланс:" + strconv.Itoa(money)
		fmt.Println(msg)
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("Fail to write HTTP response", err)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)

		msg := "Не хватает денег на проведение оплаты"
		fmt.Println(msg)
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("Fail to write HTTP response", err)
		}
	}
	mtx.Unlock()
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed) // Patch
		return
	}

	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		msg := "fail to read HTTP request body:" + err.Error()
		// msg := fmt.Sprintf("fail to read HTTP body: %v", err) ТАК ЛУЧШЕ
		fmt.Println(msg)
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("Fail to write HTTP response", err)
		}

		return
	}

	httpRequestBodyString := string(httpRequestBody)

	saveAmount, err := strconv.Atoi(httpRequestBodyString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		msg := "fail to conver HTTP body to int:" + err.Error()
		// msg := fmt.Sprintf("fail to read HTTP body: %v", err) ТАК ЛУЧШЕ
		fmt.Println(msg)
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("Fail to write HTTP response", err)
		}

		return
	}

	mtx.Lock()
	if money >= saveAmount {
		money -= saveAmount
		bank += saveAmount

		msg := "Новое значение переменной money: " + strconv.Itoa(money) + "\n" +
			"Новое значение переменной bank: " + strconv.Itoa(bank)
		fmt.Println(msg)
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("Fail to write HTTP response", err)
		}

	} else {
		w.WriteHeader(http.StatusBadRequest)

		msg := "Не хватает денег для того, чтобы положить в копилку"
		fmt.Println(msg)
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("Fail to write HTTP response", err)
		}
	}
	mtx.Unlock()

}

func main() {
	fmt.Println("Сервер запущен")

	http.HandleFunc("/pay", payHandler)
	http.HandleFunc("/save", saveHandler)

	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		fmt.Println("HTTP server error: ", err)
	}

	fmt.Println("Сервер закрыт")
}
