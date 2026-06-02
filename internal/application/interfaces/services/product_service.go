package services

import (
	"context"

	"github.com/OmarAraby/go-ecommerce/internal/domain/entities"
)

type ProductService interface {
	GetByID(ctx context.Context, id int64) (*entities.Product, error)
	List(ctx context.Context) ([]*entities.Product, error)
	Create(ctx context.Context, p *entities.Product) (*entities.Product, error)
	Update(ctx context.Context, p *entities.Product) (*entities.Product, error)
	Delete(ctx context.Context, id int64) error
}
