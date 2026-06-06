package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	orderapp "github.com/OmarAraby/go-ecommerce/internal/application/services/order"
	"github.com/OmarAraby/go-ecommerce/internal/domain"
	"github.com/OmarAraby/go-ecommerce/internal/domain/entities"
	"github.com/OmarAraby/go-ecommerce/internal/infrastructure/db"
)

var _ orderapp.OrderRepository = (*OrderRepo)(nil)

type OrderRepo struct {
	q    *db.Queries
	pool *pgxpool.Pool
}

func NewOrderRepo(pool *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{q: db.New(pool), pool: pool}
}

// Create atomically locks product rows, validates stock, decrements, and inserts order + items.
func (r *OrderRepo) Create(ctx context.Context, order *entities.Order, items []*entities.OrderItem) (*entities.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) // no-op if committed

	qtx := r.q.WithTx(tx)

	// 1. lock product rows and verify stock
	for _, item := range items {
		var stock int32
		err := tx.QueryRow(ctx,
			"SELECT stock FROM products WHERE id = $1 FOR UPDATE",
			item.ProductID,
		).Scan(&stock)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("product %d: %w", item.ProductID, domain.ErrNotFound)
			}
			return nil, fmt.Errorf("lock product %d: %w", item.ProductID, err)
		}
		if int(stock) < item.Quantity {
			return nil, fmt.Errorf("product %d: %w", item.ProductID, domain.ErrInsufficientStock)
		}
	}

	// 2. decrement stock for each item
	for _, item := range items {
		_, err := tx.Exec(ctx,
			"UPDATE products SET stock = stock - $1, updated_at = NOW() WHERE id = $2",
			item.Quantity, item.ProductID,
		)
		if err != nil {
			return nil, fmt.Errorf("decrement stock product %d: %w", item.ProductID, err)
		}
	}

	// 3. insert order
	dbOrder, err := qtx.CreateOrder(ctx, db.CreateOrderParams{
		UserID:      order.UserID,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
	})
	if err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	// 4. insert order items
	for _, item := range items {
		_, err := qtx.CreateOrderItem(ctx, db.CreateOrderItemParams{
			OrderID:     dbOrder.ID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    int32(item.Quantity),
			UnitPrice:   item.UnitPrice,
		})
		if err != nil {
			return nil, fmt.Errorf("create order item: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}
	return toOrderDomain(dbOrder), nil
}

func (r *OrderRepo) GetByID(ctx context.Context, id int64) (*entities.Order, error) {
	row, err := r.q.GetOrder(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("GetOrder: %w", err)
	}
	return toOrderDomain(row), nil
}

func (r *OrderRepo) GetUserOrder(ctx context.Context, userID, orderID int64) (*entities.Order, error) {
	row, err := r.q.GetUserOrder(ctx, db.GetUserOrderParams{ID: orderID, UserID: userID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("GetUserOrder: %w", err)
	}
	return toOrderDomain(row), nil
}

func (r *OrderRepo) GetItems(ctx context.Context, orderID int64) ([]*entities.OrderItem, error) {
	rows, err := r.q.GetOrderItems(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("GetOrderItems: %w", err)
	}
	items := make([]*entities.OrderItem, len(rows))
	for i, row := range rows {
		items[i] = toOrderItemDomain(row)
	}
	return items, nil
}

func (r *OrderRepo) GetItemsByOrderIDs(ctx context.Context, orderIDs []int64) ([]*entities.OrderItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, order_id, product_id, product_name, quantity, unit_price, created_at
		 FROM order_items WHERE order_id = ANY($1) ORDER BY order_id, id`,
		orderIDs,
	)
	if err != nil {
		return nil, fmt.Errorf("GetItemsByOrderIDs: %w", err)
	}
	defer rows.Close()

	var items []*entities.OrderItem
	for rows.Next() {
		var item entities.OrderItem
		var qty int32
		var createdAt pgtype.Timestamptz
		if err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID,
			&item.ProductName, &qty, &item.UnitPrice, &createdAt,
		); err != nil {
			return nil, fmt.Errorf("scan order item: %w", err)
		}
		item.Quantity = int(qty)
		item.CreatedAt = createdAt.Time
		items = append(items, &item)
	}
	return items, rows.Err()
}

func (r *OrderRepo) ListByUser(ctx context.Context, userID int64, params orderapp.ListOrdersParamsDTO) ([]*entities.Order, int, error) {
	var total int64
	if err := r.pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM orders WHERE user_id = $1", userID,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count orders: %w", err)
	}

	offset := int32((params.Page - 1) * params.Limit)
	rows, err := r.q.ListUserOrders(ctx, db.ListUserOrdersParams{
		UserID: userID,
		Limit:  int32(params.Limit),
		Offset: offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("ListUserOrders: %w", err)
	}

	orders := make([]*entities.Order, len(rows))
	for i, row := range rows {
		orders[i] = toOrderDomain(row)
	}
	return orders, int(total), nil
}

func toOrderDomain(o db.Order) *entities.Order {
	return &entities.Order{
		ID:          o.ID,
		UserID:      o.UserID,
		Status:      o.Status,
		TotalAmount: o.TotalAmount,
		CreatedAt:   o.CreatedAt.Time,
		UpdatedAt:   o.UpdatedAt.Time,
	}
}

func toOrderItemDomain(i db.OrderItem) *entities.OrderItem {
	return &entities.OrderItem{
		ID:          i.ID,
		OrderID:     i.OrderID,
		ProductID:   i.ProductID,
		ProductName: i.ProductName,
		Quantity:    int(i.Quantity),
		UnitPrice:   i.UnitPrice,
		CreatedAt:   i.CreatedAt.Time,
	}
}
