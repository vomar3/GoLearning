package order

import (
	"math"
	"time"
	"unicode"

	"github.com/google/uuid"
)

var epsilon64 = math.Nextafter(1.0, 2.0) - 1.0

const (
	StatusPending   string = "pending"
	StatusCooking   string = "cooking"
	StatusReady     string = "ready"
	StatusDelivered string = "delivered"
	StatusCancelled string = "cancelled"
)

// Validate string: letters only
func ValidateString(str string) bool {
	if str == "" {
		return false
	}

	for _, symbol := range str {
		if !unicode.IsLetter(symbol) && !unicode.IsSpace(symbol) {
			return false
		}
	}

	return true
}

type OrderItem struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}

type Order struct {
	ID      string      `json:"id"`
	Client  string      `json:"client"`
	Items   []OrderItem `json:"items"`
	Total   float64     `json:"total"`
	Status  string      `json:"status"` // "pending", "cooking", "ready", "delivered", "cancelled"
	Created time.Time   `json:"created"`
}

type CreateOrderRequest struct {
	Client string      `json:"client"`
	Items  []OrderItem `json:"items"`
}

type StatusUpdateRequest struct {
	Status string `json:"status"`
}

type StatsResponse struct {
	TotalOrders  int            `json:"total_orders"`
	TotalSum     float64        `json:"total_sum"`
	AverageCheck float64        `json:"average_check"`
	StatusCounts map[string]int `json:"status_counts"`
}

// Check data for a separate dish
func (order OrderItem) Validate() bool {
	if !ValidateString(order.Name) {
		return false
	}

	if order.Price <= 0 || order.Qty <= 0 {
		return false
	}

	return true
}

// Validate status of order
func (status StatusUpdateRequest) Validate() bool {
	switch status.Status {
	case StatusPending, StatusCooking, StatusReady, StatusDelivered, StatusCancelled:
		return true
	default:
		return false
	}
}

// Validate all order without ID (RN)
func (order *Order) Validate() bool {
	if order == nil {
		return false
	}

	if !ValidateString(order.Client) {
		return false
	}

	if order.Total < epsilon64 {
		return false
	}

	switch order.Status {
	case StatusPending, StatusCooking, StatusReady, StatusDelivered, StatusCancelled:
	default:
		return false
	}

	for _, val := range order.Items {
		if !val.Validate() {
			return false
		}
	}

	return true
}

func (newOrder *CreateOrderRequest) Validate() bool {
	if newOrder == nil {
		return false
	}

	if !ValidateString(newOrder.Client) {
		return false
	}

	for _, val := range newOrder.Items {
		if !val.Validate() {
			return false
		}
	}

	return true
}

// Calculate price from 1 order Item
func (item OrderItem) CalculatePositionPrice() float64 {
	return item.Price * float64(item.Qty)
}

// Calculate all order price
func (order *CreateOrderRequest) CalculateOrderPrice() float64 {
	var price float64
	for _, data := range order.Items {
		price += data.CalculatePositionPrice()
	}

	return price
}

/*func (order *Order) CalculateOrderPrice() float64 {
	var price float64
	for _, data := range order.Items {
		price += data.CalculatePositionPrice()
	}

	return price
}*/

func CreateOrder(orderData CreateOrderRequest) *Order {
	newUUID := uuid.New()

	order := &Order{
		ID:      newUUID.String(),
		Client:  orderData.Client,
		Items:   orderData.Items,
		Total:   orderData.CalculateOrderPrice(),
		Status:  StatusPending,
		Created: time.Now(),
	}

	return order
}

func ChangeStatus(oldStatus, newStatus string) bool {
	switch oldStatus {
	case StatusPending:
		return newStatus == StatusCancelled || newStatus == StatusCooking
	case StatusCooking:
		return newStatus == StatusReady
	case StatusReady:
		return newStatus == StatusDelivered
	case StatusDelivered, StatusCancelled:
		return false
	}

	return false

}

func ActiveOrders(str string) bool {
	switch str {
	case StatusPending, StatusCooking, StatusReady:
		return true
	}

	return false
}

func CreateStats() *StatsResponse {
	return &StatsResponse{
		TotalOrders:  0,
		TotalSum:     0,
		AverageCheck: 0,
		StatusCounts: map[string]int{
			StatusPending:   0,
			StatusCooking:   0,
			StatusReady:     0,
			StatusDelivered: 0,
			StatusCancelled: 0,
		},
	}
}
