package models

import (
	"fmt"
	"time"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusProcessed  OrderStatus = "processed"
	OrderStatusFailed     OrderStatus = "failed"
)

type OrderRequest struct {
	ID    string `json:"id"`
	Item  string `json:"item"`
	Price int    `json:"price"`
}

type Order struct {
	ID        string      `json:"id"`
	Item      string      `json:"item"`
	Price     int         `json:"price"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type OrderResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg,omitempty"`
}

func (o *OrderRequest) Validate() error {
	if o.ID == "" {
		return fmt.Errorf("field 'id' is required")
	}
	if o.Item == "" {
		return fmt.Errorf("field 'item' is required")
	}
	if o.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	return nil
}
