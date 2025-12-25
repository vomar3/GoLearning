package main

import (
	"ManagementAPI/order" // из go.mod название
	"ManagementAPI/storage"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
)

// Parsing data from JSON. Make new order with status "pending"
func HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed) // Need Post
		_, _ = w.Write([]byte("Bad status, use POST"))
		return
	}

	var NewOrder order.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&NewOrder); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error with parse JSON"))
		return
	}

	if !NewOrder.Validate() {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Bad Request %+v", NewOrder)
		_, _ = w.Write([]byte(msg))
		return
	}

	order := order.CreateOrder(NewOrder)
	storage.Mtx.Lock()
	defer storage.Mtx.Unlock()

	storage.Orders[order.ID] = order

	w.WriteHeader(http.StatusOK)
	msg := fmt.Sprintf("New order with ID: %s", order.ID)
	_, _ = w.Write([]byte(msg))
}

func HandleCheckOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed) // need GET
		_, _ = w.Write([]byte("Bad status, use GET"))
		return
	}

	id := r.PathValue("id")
	if err := uuid.Validate(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Error with id's format: %s. Please check id.", id)
		_, _ = w.Write([]byte(msg))
		return
	}

	storage.Mtx.RLock()
	defer storage.Mtx.RUnlock()

	val, ok := storage.Orders[id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Error the specific id: %s, does not exist", id)
		_, _ = w.Write([]byte(msg))
		return
	}

	request, err := json.Marshal(val)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error with convert request to JSON"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(request)
}

func HandleChangeStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.Header().Set("Allow", "PATCH")
		w.WriteHeader(http.StatusMethodNotAllowed) // need PATCH
		_, _ = w.Write([]byte("Bad status, use PATCH"))
		return
	}

	id := r.PathValue("id")
	if err := uuid.Validate(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Error with id's format: %s. Please check id.", id)
		_, _ = w.Write([]byte(msg))
		return
	}

	storage.Mtx.Lock()
	defer storage.Mtx.Unlock()

	ord, ok := storage.Orders[id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Error the specific id: %s, does not exist", id)
		_, _ = w.Write([]byte(msg))
		return
	}

	var status order.StatusUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error with convert from JSON"))
		return
	}

	if !status.Validate() {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Bad status: %s", status)
		_, _ = w.Write([]byte(msg))
		return
	}

	if !order.ChangeStatus(ord.Status, status.Status) {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Can't use new status (%s) right now. You have status: %s", status.Status, storage.Orders[id].Status)
		_, _ = w.Write([]byte(msg))
		return
	}

	ord.Status = status.Status

	w.WriteHeader(http.StatusOK)
	msg := fmt.Sprintf("Succesfull. New status is %s", status.Status)
	_, _ = w.Write([]byte(msg))
}

func HandleActiveOrders(w http.ResponseWriter, r *http.Request) {
	activeOrders := make([]*order.Order, 0)

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Bad status, use GET"))
		return
	}

	queryStatus := r.URL.Query().Get("status")
	if !order.ActiveOrders(queryStatus) {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Bad status: %s. Need active status ('pending', 'ready', 'cooking')", queryStatus)
		_, _ = w.Write([]byte(msg))
		return
	}

	storage.Mtx.RLock()
	defer storage.Mtx.RUnlock()

	for _, value := range storage.Orders {
		if value.Status == queryStatus {
			activeOrders = append(activeOrders, value)
		}
	}

	response, err := json.Marshal(activeOrders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error with convert request to JSON"))
		return
	}

	w.Header().Set("Content-Type", "application/json") // Пояснительная надпись о том, что пользователю улетит JSON, можно использовать с методами JSON.
	w.WriteHeader(http.StatusOK)
	w.Write(response) // json.Marshal возвращает []byte, поэтому нет каста
}

func HandleCancelOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Bad status, use POST"))
		return
	}

	id := r.PathValue("id")
	if err := uuid.Validate(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Error with id's format: %s. Please check id.", id)
		_, _ = w.Write([]byte(msg))
		return
	}

	storage.Mtx.Lock()
	defer storage.Mtx.Unlock()

	myOrder, ok := storage.Orders[id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Error the specific id: %s, does not exist", id)
		_, _ = w.Write([]byte(msg))
		return
	}

	if myOrder.Status != order.StatusPending {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("You can't cancel your order, because your status is %s", myOrder.Status)
		_, _ = w.Write([]byte(msg))
		return
	}

	myOrder.Status = order.StatusCancelled
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("You have cancelled your order"))
}

func HandleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Bad status, use GET"))
		return
	}

	stats := order.CreateStats()

	storage.Mtx.RLock()
	for _, value := range storage.Orders {
		stats.TotalOrders++
		stats.StatusCounts[value.Status]++
		if value.Status != order.StatusCancelled {
			stats.TotalSum += value.Total
		}
	}
	storage.Mtx.RUnlock()

	if stats.TotalOrders > 0 {
		stats.AverageCheck = stats.TotalSum / float64(stats.TotalOrders)
	} else {
		stats.AverageCheck = 0
	}

	response, err := json.Marshal(stats)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error with convert request to JSON"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {
	// Создаем не как http.ListenAndServe, т.к. в том случае таймауты бесконечные
	// Тут мы сами настраиваем таймауты
	srv := &http.Server{
		Addr: ":9091",
		// Если клиент подключился, но молчит 100 сек — разрываем
		IdleTimeout: 100 * time.Second,
		// Если клиент слишком долго шлет запрос (медленный инет)
		ReadTimeout: 10 * time.Second,
		// Если мы (сервер) тупим и долго формируем ответ
		WriteTimeout: 10 * time.Second,
	}

	http.HandleFunc("/order/create", HandleCreate)
	http.HandleFunc("/order/{id}", HandleCheckOrder)
	http.HandleFunc("/order/{id}/status", HandleChangeStatus)
	http.HandleFunc("/order/active", HandleActiveOrders)
	http.HandleFunc("/order/{id}/cancel", HandleCancelOrder)
	http.HandleFunc("/stats", HandleStats)

	// В горутине, чтобы не блокался мейн
	go func() {
		fmt.Println("The server is starting on :9091")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Listen error: %s\n", err)
		}
	}()

	// Создаем канал, чтобы охватывать ошибку, и как только получается, то завершаем работу сервера
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("\nShutting down server...")

	// даем 5 секунд на выполнение оставшихся задач
	// закончились раньше = сразу закрывается сервер
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Закрываем сервак
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}

	fmt.Println("Server exiting")
}
