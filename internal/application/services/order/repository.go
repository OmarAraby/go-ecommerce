package order

import (
	"context"

	"github.com/OmarAraby/go-ecommerce/internal/domain/entities"
)

// OrderRepository handles order persistence and the stock-decrement transaction.
type OrderRepository interface {
	// Create atomically: locks product rows, checks stock, decrements, inserts order + items.
	Create(ctx context.Context, order *entities.Order, items []*entities.OrderItem) (*entities.Order, error)
	GetByID(ctx context.Context, id int64) (*entities.Order, error)
	GetUserOrder(ctx context.Context, userID, orderID int64) (*entities.Order, error)
	GetItems(ctx context.Context, orderID int64) ([]*entities.OrderItem, error)
	GetItemsByOrderIDs(ctx context.Context, orderIDs []int64) ([]*entities.OrderItem, error)
	ListByUser(ctx context.Context, userID int64, params ListOrdersParamsDTO) ([]*entities.Order, int, error)
}

// ProductLookup is a minimal interface for reading product data needed by the order service.
// Defined here (at the consumer) — the product.Repository satisfies this automatically.
type ProductLookup interface {
	GetByID(ctx context.Context, id int64) (*entities.Product, error)
}
