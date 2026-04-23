package models

import "testing"

func TestOrderRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		order   OrderRequest
		wantErr bool
	}{
		{
			name:    "valid order",
			order:   OrderRequest{ID: "order-1", Item: "book", Price: 1200},
			wantErr: false,
		},
		{
			name:    "missing id",
			order:   OrderRequest{Item: "book", Price: 1200},
			wantErr: true,
		},
		{
			name:    "missing item",
			order:   OrderRequest{ID: "order-1", Price: 1200},
			wantErr: true,
		},
		{
			name:    "non positive price",
			order:   OrderRequest{ID: "order-1", Item: "book"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.order.Validate()
			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		})
	}
}
