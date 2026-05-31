package repositories

import (
	"context"

	"github.com/OmarAraby/go-ecommerce/internal/domain"
)

type ProductRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	List(ctx context.Context) ([]*domain.Product, error)
	Create(ctx context.Context, p *domain.Product) (*domain.Product, error)
	Update(ctx context.Context, p *domain.Product) (*domain.Product, error)
	Delete(ctx context.Context, id int64) error
}
