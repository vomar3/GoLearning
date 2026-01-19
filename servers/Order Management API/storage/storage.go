package storage

import (
	"ManagementAPI/order"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
)

type PostgresStorage struct {
	db *pgx.Conn
}

type Storage interface {
	AddOrder(ord *order.Order) error
	Get(id string) (*order.Order, error)
	UpdateStatus(id, status string) error
	CancelOrder(id string) error
	GetByStatus(status string) ([]order.Order, error)
	GetAllStats() (*order.StatsResponse, error)
}

func (s *PostgresStorage) Get(id string) (*order.Order, error) {
	var ord order.Order
	var itemsJSON []byte

	query := `
		SELECT id, client, status, total, items, created_at
		FROM orders
		WHERE id = @id
	`

	args := pgx.NamedArgs{"id": id}

	err := s.db.QueryRow(context.Background(), query, args).Scan(
		&ord.ID,
		&ord.Client,
		&ord.Status,
		&ord.Total,
		&itemsJSON,
		&ord.Created,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("Get: order not found: %s", id)
	}

	if err != nil {
		return nil, fmt.Errorf("Get: failed to get order: %w", err)
	}

	if err := json.Unmarshal(itemsJSON, &ord.Items); err != nil {
		return nil, fmt.Errorf("Get: failed to unmarshal items: %w", err)
	}

	return &ord, nil
}

func (s *PostgresStorage) UpdateStatus(id, status string) error {
	ord, err := s.Get(id)
	if err != nil {
		return fmt.Errorf("UpdateStatus: %w", err)
	}

	if ok := order.ChangeStatus(ord.Status, status); !ok {
		return fmt.Errorf("UpdateStatus: You couldn't update status")
	}

	query := `
		UPDATE orders
		SET status = @status
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id":     id,
		"status": status,
	}

	_, err = s.db.Exec(context.Background(), query, args)

	if err != nil {
		return fmt.Errorf("UpdateStatus: %w", err)
	}

	return nil
}

func (s *PostgresStorage) CancelOrder(id string) error {
	ord, err := s.Get(id)
	if err != nil {
		return fmt.Errorf("CancelOrder: %w", err)
	}

	if !(ord.Status == order.StatusPending) {
		return fmt.Errorf("CancelOrder: you can't cancel your order because your status is %s", ord.Status)
	}

	query := `
		UPDATE orders
		SET status = @status
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id":     id,
		"status": order.StatusCancelled,
	}

	_, err = s.db.Exec(context.Background(), query, args)
	if err != nil {
		return fmt.Errorf("CancelOrder: %w", err)
	}

	return nil
}

func (s *PostgresStorage) GetByStatus(status string) ([]order.Order, error) {
	answer := make([]order.Order, 0)

	query := `
		SELECT id, client, status, total, items, created_at
		FROM orders
		WHERE status = @status
	`

	args := pgx.NamedArgs{
		"status": status,
	}

	rows, err := s.db.Query(context.Background(), query, args)
	if err != nil {
		return nil, fmt.Errorf("GetByStatus: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var ord order.Order
		var itemsJSON []byte

		if err := rows.Scan(
			&ord.ID,
			&ord.Client,
			&ord.Status,
			&ord.Total,
			&itemsJSON,
			&ord.Created,
		); err != nil {
			return nil, fmt.Errorf("GetByStatus: %w", err)
		}

		if err := json.Unmarshal(itemsJSON, &ord.Items); err != nil {
			return nil, fmt.Errorf("GetByStatus: failed to unmarshal items: %w", err)
		}

		answer = append(answer, ord)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetByStatus: %w", err)
	}

	return answer, nil
}

func (s *PostgresStorage) GetAllStats() (*order.StatsResponse, error) {
	stats := order.CreateStats()

	queryTotal := `
        SELECT 
            COUNT(*) as total_count,
            COALESCE(SUM(total), 0) as total_sum,
            COALESCE(AVG(total), 0) as avg_check
        FROM orders
        WHERE status != @cancelledStatus
    `

	args := pgx.NamedArgs{
		"cancelledStatus": order.StatusCancelled,
	}

	err := s.db.QueryRow(context.Background(), queryTotal, args).Scan(
		&stats.TotalOrders,
		&stats.TotalSum,
		&stats.AverageCheck,
	)

	if err != nil {
		return nil, fmt.Errorf("GetAllStats: %w", err)
	}

	queryStatus := `
        SELECT status, COUNT(*)
        FROM orders
        GROUP BY status
    `

	rows, err := s.db.Query(context.Background(), queryStatus)
	if err != nil {
		return nil, fmt.Errorf("GetAllStats: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var status string
		var count int

		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("GetAllStats: %w", err)
		}

		stats.StatusCounts[status] = count
	}

	return stats, nil
}

// Код подключения
func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	ctx := context.Background()

	// 1. Подключаемся через pgx
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// 2. Проверяем связь
	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	// 3. Создаем таблицу (если нет)
	// Обрати внимание: тут обычный SQL, pgx его отлично понимает
	query := `
	CREATE TABLE IF NOT EXISTS orders (
		id TEXT PRIMARY KEY,
		client TEXT NOT NULL,
		status TEXT NOT NULL,
		total NUMERIC NOT NULL,
		items JSONB NOT NULL,
		created_at TIMESTAMP NOT NULL
	);
	`
	_, err = conn.Exec(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &PostgresStorage{db: conn}, nil
}

func (s *PostgresStorage) AddOrder(ord *order.Order) error {
	// Postgres не умеет хранить Go-слайс []OrderItem напрямую.
	// Мы превращаем его в JSON-строку (массив байт).
	// Было: [{Name: "Pizza", Price: 100}] -> Стало: `[{"name":"Pizza","price":100}]`
	itemsJSON, err := json.Marshal(ord.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}

	// Подготовка SQL-запроса
	// Используем именованные аргументы (@id, @client)
	query := `
        INSERT INTO orders (id, client, status, total, items, created_at) 
        VALUES (@id, @client, @status, @total, @items, @created_at)
    `

	// Связывание данных (Arguments)
	// Мы создаем карту (map), где ключи — это имена из запроса (@id),
	// а значения — это реальные данные из нашей структуры.
	args := pgx.NamedArgs{
		"id":         ord.ID,
		"client":     ord.Client,
		"status":     ord.Status,
		"total":      ord.Total,
		"items":      itemsJSON, // Важно! Передаем именно JSON-байты
		"created_at": ord.Created,
	}

	// Отправка (Exec)
	// context.Background() — нужен для отмены запроса.
	// query — наш текст SQL.
	// args — наши данные.
	// Функция Exec сама всё склеит безопасно и отправит в Postgres.
	_, err = s.db.Exec(context.Background(), query, args)

	// Если Postgres вернул ошибку (например, такой ID уже есть), мы её поймаем тут.
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	return nil
}

type MemoryStorage struct {
	data map[string]*order.Order
	mtx  sync.RWMutex
}

func NewStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]*order.Order),
	}
}

func (storage *MemoryStorage) AddOrder(ord *order.Order) error {
	if ord == nil {
		return fmt.Errorf("AddOrder: The pointer must not be a nil")
	}

	storage.mtx.Lock()
	defer storage.mtx.Unlock()

	storage.data[ord.ID] = ord

	return nil
}

func (storage *MemoryStorage) Get(id string) (*order.Order, error) {
	storage.mtx.RLock()
	defer storage.mtx.RUnlock()

	val, ok := storage.data[id]
	if !ok {
		return nil, fmt.Errorf("Error the specific id: %s, does not exist", id)
	}

	copy := *val

	return &copy, nil
}

func (storage *MemoryStorage) GetByStatus(status string) ([]order.Order, error) {
	activeOrders := make([]order.Order, 0)
	storage.mtx.RLock()
	defer storage.mtx.RUnlock()

	for _, value := range storage.data {
		if value.Status == status {
			activeOrders = append(activeOrders, *value)
		}
	}

	return activeOrders, nil
}

func (storage *MemoryStorage) UpdateStatus(id, newStatus string) error {
	storage.mtx.Lock()
	defer storage.mtx.Unlock()

	ord, ok := storage.data[id]
	if !ok {
		return fmt.Errorf("UpdateStatus: ID %s not found", id)
	}

	if order.ChangeStatus(ord.Status, newStatus) {
		ord.Status = newStatus
	} else {
		return fmt.Errorf("Can't use new status (%s) right now. You have status: %s", newStatus, ord.Status)
	}

	return nil
}

func (storage *MemoryStorage) CancelOrder(id string) error {
	storage.mtx.Lock()
	defer storage.mtx.Unlock()

	ord, ok := storage.data[id]
	if !ok {
		return fmt.Errorf("CancelOrder: ID %s not found", id)
	}

	if ord.Status != order.StatusPending {
		return fmt.Errorf("You can't cancel your order, because your status is %s", ord.Status)
	}

	ord.Status = order.StatusCancelled

	return nil
}

func (storage *MemoryStorage) GetAllStats() (*order.StatsResponse, error) {
	stats := order.CreateStats()

	storage.mtx.RLock()
	for _, value := range storage.data {
		stats.TotalOrders++
		stats.StatusCounts[value.Status]++
		if value.Status != order.StatusCancelled {
			stats.TotalSum += value.Total
		}
	}
	storage.mtx.RUnlock()

	if stats.TotalOrders > 0 {
		stats.AverageCheck = stats.TotalSum / float64(stats.TotalOrders)
	} else {
		stats.AverageCheck = 0
	}

	copy := *stats
	return &copy, nil
}
