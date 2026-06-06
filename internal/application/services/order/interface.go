package order

import "context"

type Service interface {
	Create(ctx context.Context, dto CreateOrderDTO) (*OrderResponseDTO, error)
	GetByID(ctx context.Context, userID, orderID int64) (*OrderResponseDTO, error)
	ListByUser(ctx context.Context, userID int64, params ListOrdersParamsDTO) ([]OrderResponseDTO, int, error)
}
