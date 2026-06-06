package order

import (
	"context"
	"fmt"

	"github.com/OmarAraby/go-ecommerce/internal/domain"
	"github.com/OmarAraby/go-ecommerce/internal/domain/entities"
)

var _ Service = (*service)(nil)

type service struct {
	orders   OrderRepository
	products ProductLookup
}

func NewService(orders OrderRepository, products ProductLookup) Service {
	return &service{orders: orders, products: products}
}

func (s *service) Create(ctx context.Context, dto CreateOrderDTO) (*OrderResponseDTO, error) {
	if len(dto.Items) == 0 {
		return nil, fmt.Errorf("order must have at least one item: %w", domain.ErrNotFound)
	}

	// 1. fetch products and build snapshots — validate they exist
	type enriched struct {
		item    CreateOrderItemDTO
		product *entities.Product
	}
	enrich := make([]enriched, len(dto.Items))
	for i, item := range dto.Items {
		p, err := s.products.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %d: %w", item.ProductID, domain.ErrNotFound)
		}
		enrich[i] = enriched{item: item, product: p}
	}

	// 2. calculate total
	var total float64
	for _, e := range enrich {
		total += e.product.Price * float64(e.item.Quantity)
	}

	// 3. build domain objects (price/name snapshot captured here)
	order := &entities.Order{
		UserID:      dto.UserID,
		Status:      "pending",
		TotalAmount: total,
	}
	items := make([]*entities.OrderItem, len(enrich))
	for i, e := range enrich {
		items[i] = &entities.OrderItem{
			ProductID:   e.product.ID,
			ProductName: e.product.Name,
			Quantity:    e.item.Quantity,
			UnitPrice:   e.product.Price,
		}
	}

	// 4. persist atomically — repo handles stock lock + decrement + insert
	created, err := s.orders.Create(ctx, order, items)
	if err != nil {
		return nil, err
	}

	// 5. fetch items to build response (repo inserted them)
	orderItems, err := s.orders.GetItems(ctx, created.ID)
	if err != nil {
		return nil, err
	}

	r := toResponse(created, orderItems)
	return &r, nil
}

func (s *service) GetByID(ctx context.Context, userID, orderID int64) (*OrderResponseDTO, error) {
	order, err := s.orders.GetUserOrder(ctx, userID, orderID)
	if err != nil {
		return nil, err
	}
	items, err := s.orders.GetItems(ctx, orderID)
	if err != nil {
		return nil, err
	}
	r := toResponse(order, items)
	return &r, nil
}

func (s *service) ListByUser(ctx context.Context, userID int64, params ListOrdersParamsDTO) ([]OrderResponseDTO, int, error) {
	orders, total, err := s.orders.ListByUser(ctx, userID, params)
	if err != nil {
		return nil, 0, err
	}
	if len(orders) == 0 {
		return []OrderResponseDTO{}, total, nil
	}

	// batch-fetch all items in one query
	ids := make([]int64, len(orders))
	for i, o := range orders {
		ids[i] = o.ID
	}
	allItems, err := s.orders.GetItemsByOrderIDs(ctx, ids)
	if err != nil {
		return nil, 0, err
	}
	itemMap := make(map[int64][]*entities.OrderItem)
	for _, item := range allItems {
		itemMap[item.OrderID] = append(itemMap[item.OrderID], item)
	}

	result := make([]OrderResponseDTO, len(orders))
	for i, o := range orders {
		result[i] = toResponse(o, itemMap[o.ID])
	}
	return result, total, nil
}

func toResponse(o *entities.Order, items []*entities.OrderItem) OrderResponseDTO {
	itemDTOs := make([]OrderItemResponseDTO, len(items))
	for i, item := range items {
		itemDTOs[i] = OrderItemResponseDTO{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Subtotal:    float64(item.Quantity) * item.UnitPrice,
		}
	}
	return OrderResponseDTO{
		ID:          o.ID,
		Status:      o.Status,
		TotalAmount: o.TotalAmount,
		Items:       itemDTOs,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}
