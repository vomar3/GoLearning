package models

import "fmt"

type OrderRequest struct {
	ID    string `json:"id"`
	Item  string `json:"item"`
	Price int    `json:"price"`
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
