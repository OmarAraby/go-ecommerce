package product

import (
	"context"

	"github.com/OmarAraby/go-ecommerce/internal/domain/entities"
)

// Repository is the contract the infrastructure layer must satisfy.
// Defined here (at the consumer) — not at the implementation.
type Repository interface {
	GetByID(ctx context.Context, id int64) (*entities.Product, error)
	List(ctx context.Context, params ListParams) ([]*entities.Product, int, error)
	Create(ctx context.Context, p *entities.Product) (*entities.Product, error)
	Update(ctx context.Context, p *entities.Product) (*entities.Product, error)
	Delete(ctx context.Context, id int64) error

	// Image operations
	GetImages(ctx context.Context, productID int64) ([]*entities.ProductImage, error)
	GetImagesByProductIDs(ctx context.Context, productIDs []int64) ([]*entities.ProductImage, error)
	CountImages(ctx context.Context, productID int64) (int, error)
	AddImage(ctx context.Context, productID int64, url string, isMain bool) (*entities.ProductImage, error)
	DeleteImage(ctx context.Context, productID, imageID int64) error
	SetMainImage(ctx context.Context, productID, imageID int64) error
	GetImage(ctx context.Context, productID, imageID int64) (*entities.ProductImage, error)
}
