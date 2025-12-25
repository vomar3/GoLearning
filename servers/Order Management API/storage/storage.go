package storage

import (
	"ManagementAPI/order"
	"sync"
)

var (
	// Хранилище: ID заказа -> Структура заказа
	Orders = make(map[string]*order.Order)

	// Мьютекс: Обязателен, так как map не потокобезопасна
	// (одновременно могут прийти запрос на чтение и на создание)
	Mtx sync.RWMutex
)
