package application

import (
	"context"

	"github.com/OmarAraby/go-ecommerce/internal/application/interfaces/repositories"
	"github.com/OmarAraby/go-ecommerce/internal/application/interfaces/services"
	"github.com/OmarAraby/go-ecommerce/internal/domain/entities"
)

var _ services.ProductService = (*ProductService)(nil)

type ProductService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetByID(ctx context.Context, id int64) (*entities.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) List(ctx context.Context) ([]*entities.Product, error) {
	return s.repo.List(ctx)
}

func (s *ProductService) Create(ctx context.Context, p *entities.Product) (*entities.Product, error) {
	return s.repo.Create(ctx, p)
}

func (s *ProductService) Update(ctx context.Context, p *entities.Product) (*entities.Product, error) {
	return s.repo.Update(ctx, p)
}

func (s *ProductService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
