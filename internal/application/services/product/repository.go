package product

import (
	"context"

	"github.com/OmarAraby/go-ecommerce/internal/domain/entities"
)

// Repository is the contract the infrastructure layer must satisfy.
// Defined here (at the consumer) — not at the implementation.
type Repository interface {
	GetByID(ctx context.Context, id int64) (*entities.Product, error)
	List(ctx context.Context) ([]*entities.Product, error)
	Create(ctx context.Context, p *entities.Product) (*entities.Product, error)
	Update(ctx context.Context, p *entities.Product) (*entities.Product, error)
	Delete(ctx context.Context, id int64) error
}
