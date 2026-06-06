package order

import "time"

type CreateOrderItemDTO struct {
	ProductID int64 `json:"product_id" validate:"required,gt=0"`
	Quantity  int   `json:"quantity"   validate:"required,gte=1"`
}

type CreateOrderDTO struct {
	UserID int64
	Items  []CreateOrderItemDTO
}

// ──────────────────────────────────────────────

type OrderItemResponseDTO struct {
	ID          int64   `json:"id"`
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Subtotal    float64 `json:"subtotal"` // quantity × unit_price
}

type OrderResponseDTO struct {
	ID          int64                  `json:"id"`
	Status      string                 `json:"status"`
	TotalAmount float64                `json:"total_amount"`
	Items       []OrderItemResponseDTO `json:"items"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type ListOrdersParamsDTO struct {
	Page  int
	Limit int
}
