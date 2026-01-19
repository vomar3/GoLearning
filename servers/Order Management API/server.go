package main

import (
	"ManagementAPI/order" // из go.mod название
	"ManagementAPI/storage"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type Server struct {
	store  storage.Storage
	logger *slog.Logger
	// info = просто какая-то информация, мы создали, отменили, измениили и т.д.
	// warn = ошибка со стороны клиента: прислал не то. Сервер работает корректно
	// error = ошибка на стороне сервака: не смогли в json, что-то упало
}

func NewServer(store storage.Storage, logger *slog.Logger) *Server {
	return &Server{
		store:  store,
		logger: logger,
	}
}

// Parsing data from JSON. Make new order with status "pending"
func (s *Server) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed) // Need Post
		_, _ = w.Write([]byte("Bad status, use POST"))
		return
	}

	var NewOrder order.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&NewOrder); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("failed to decode json", "error", err, "method", "HandleCreate")
		_, _ = w.Write([]byte("Error with parse JSON"))
		return
	}

	if !NewOrder.Validate() {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Bad Request %+v", NewOrder)
		s.logger.Warn("invalid input data", "error", msg, "method", "HandleCreate")
		_, _ = w.Write([]byte(msg))
		return
	}

	order := order.CreateOrder(NewOrder)

	if err := s.store.AddOrder(order); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("failed to save order", "order_id", order.ID, "error", err, "method", "HandleCreate")
		msg := err.Error()
		_, _ = w.Write([]byte(msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	s.logger.Info("order created", "id", order.ID, "client", order.Client, "total", order.Total, "method", "HandleCreate")
	msg := fmt.Sprintf("New order with ID: %s", order.ID)
	_, _ = w.Write([]byte(msg))
}

func (s *Server) HandleCheckOrder(w http.ResponseWriter, r *http.Request) {
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
		s.logger.Warn("invalid input id", "id", id, "error", msg, "method", "HandleCheckOrder")
		_, _ = w.Write([]byte(msg))
		return
	}

	ord, err := s.store.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := err.Error()
		s.logger.Warn("order not found", "id", id, "error", msg, "method", "HandleCheckOrder")
		_, _ = w.Write([]byte(msg))
		return
	}

	// s.store.Get returns a copy, so json.Marshal is safe without lock.
	request, err := json.Marshal(ord)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := "Error with convert request to JSON"
		s.logger.Error("failed to make JSON", "error", msg, "method", "HandleCheckOrder")
		_, _ = w.Write([]byte(msg))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(request)
	// Не пишется логгирование в гете, т.к. станет шумом в логах
}

func (s *Server) HandleChangeStatus(w http.ResponseWriter, r *http.Request) {
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
		s.logger.Warn("indalid input id", "id", id, "error", msg, "method", "HandleChangeStatus")
		_, _ = w.Write([]byte(msg))
		return
	}

	var status order.StatusUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := "Error with convert from JSON"
		s.logger.Error("failed to decode JSON", "error", msg, "method", "HandleChangeStatus")
		_, _ = w.Write([]byte(msg))
		return
	}

	if !status.Validate() {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Bad status: %s", status)
		s.logger.Warn("invalid input data", "error", msg, "method", "HandleChangeStatus")
		_, _ = w.Write([]byte(msg))
		return
	}

	if err := s.store.UpdateStatus(id, status.Status); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := err.Error()
		s.logger.Warn("invalid input data", "error", msg, "method", "HandleChangeStatus")
		_, _ = w.Write([]byte(msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	msg := fmt.Sprintf("Succesfull. New status is %s", status.Status)
	s.logger.Info("status changed", "msg", msg, "method", "HandleChangeStatus")
	_, _ = w.Write([]byte(msg))
}

func (s *Server) HandleActiveOrders(w http.ResponseWriter, r *http.Request) {
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
		s.logger.Warn("invalid input status", "status", queryStatus, "error", msg, "method", "HandleActiveOrders")
		_, _ = w.Write([]byte(msg))
		return
	}

	activeOrders, err := s.store.GetByStatus(queryStatus)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := "Error with get orders by status"
		s.logger.Error("Error with get", "error", msg, "method", "HandleActiveOrders")
		_, _ = w.Write([]byte(msg))
	}

	response, err := json.Marshal(activeOrders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := "Error with convert request to JSON"
		s.logger.Error("failed to code to JSON", "error", msg, "method", "HandleActiveOrders")
		_, _ = w.Write([]byte(msg))
		return
	}

	w.Header().Set("Content-Type", "application/json") // Пояснительная надпись о том, что пользователю улетит JSON, можно использовать с методами JSON.
	w.WriteHeader(http.StatusOK)
	w.Write(response) // json.Marshal возвращает []byte, поэтому нет каста
}

func (s *Server) HandleCancelOrder(w http.ResponseWriter, r *http.Request) {
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
		s.logger.Warn("invalid input id", "id", id, "error", msg, "method", "HandleCancelOrder")
		_, _ = w.Write([]byte(msg))
		return
	}

	err := s.store.CancelOrder(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := err.Error()
		s.logger.Warn("invalid input data", "error", msg, "method", "HandleCancelOrder")
		_, _ = w.Write([]byte(msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	msg := fmt.Sprintf("You have cancelled your order, id : %s", id)
	s.logger.Info("order canceled", "order_id", id, "method", "HandleCancelOrder")
	_, _ = w.Write([]byte(msg))
}

func (s *Server) HandleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Bad status, use GET"))
		return
	}

	stats, err := s.store.GetAllStats()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := "couldn't collect the data"
		s.logger.Error(msg, "error", err, "method", "HandleStats")
		_, _ = w.Write([]byte(msg))
	}

	response, err := json.Marshal(stats)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := "Error with convert request to JSON"
		s.logger.Error("Failed to code to JSON", "error", msg, "method", "HandleStats")
		_, _ = w.Write([]byte(msg))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	connStr := "postgres://postgres:1234@localhost:5432/postgres"

	store, err := storage.NewPostgresStorage(connStr)
	if err != nil {
		// Если не смогли подключиться — нет смысла запускать сервер, падаем сразу
		logger.Error("failed to connect to db", "error", err)
		os.Exit(1)
	}

	myServer := NewServer(store, logger)

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

	http.HandleFunc("/order/create", myServer.HandleCreate)
	http.HandleFunc("/order/{id}", myServer.HandleCheckOrder)
	http.HandleFunc("/order/{id}/status", myServer.HandleChangeStatus)
	http.HandleFunc("/order/active", myServer.HandleActiveOrders)
	http.HandleFunc("/order/{id}/cancel", myServer.HandleCancelOrder)
	http.HandleFunc("/stats", myServer.HandleStats)

	// В горутине, чтобы не блокался мейн
	go func() {
		logger.Info("server starting", "addr", ":9091")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Info("Listen error:", "error", err)
		}
	}()

	// Создаем канал, чтобы охватывать ошибку, и как только получается, то завершаем работу сервера
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// даем 5 секунд на выполнение оставшихся задач
	// закончились раньше = сразу закрывается сервер
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Закрываем сервак
	if err := srv.Shutdown(ctx); err != nil {
		logger.Info("Server forced to shutdown:", "error", err)
	}

	logger.Info("Server exiting")
}
