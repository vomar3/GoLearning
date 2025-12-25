package storage

import (
	"ManagementAPI/order"
	"fmt"
	"sync"
)

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

func (storage *MemoryStorage) GetByStatus(status string) []order.Order {
	activeOrders := make([]order.Order, 0)
	storage.mtx.RLock()
	defer storage.mtx.RUnlock()

	for _, value := range storage.data {
		if value.Status == status {
			activeOrders = append(activeOrders, *value)
		}
	}

	return activeOrders
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

func (storage *MemoryStorage) GetAllStats() *order.StatsResponse {
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
	return &copy
}
