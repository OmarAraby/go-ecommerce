package entities

import "time"

type Order struct {
	ID          int64
	UserID      int64
	Status      string
	TotalAmount float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OrderItem struct {
	ID          int64
	OrderID     int64
	ProductID   int64
	ProductName string // snapshot — preserved even if product is renamed/deleted
	Quantity    int
	UnitPrice   float64 // snapshot — preserved even if product price changes
	CreatedAt   time.Time
}
