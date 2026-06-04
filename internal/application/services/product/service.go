package product

import (
	"context"

	"github.com/OmarAraby/go-ecommerce/internal/domain/entities"
)

var _ Service = (*service)(nil)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetByID(ctx context.Context, id int64) (*entities.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) List(ctx context.Context) ([]*entities.Product, error) {
	return s.repo.List(ctx)
}

func (s *service) Create(ctx context.Context, dto CreateProductDTO) (*entities.Product, error) {
	return s.repo.Create(ctx, &entities.Product{
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		Stock:       dto.Stock,
	})
}

func (s *service) Update(ctx context.Context, dto UpdateProductDTO) (*entities.Product, error) {
	return s.repo.Update(ctx, &entities.Product{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		Stock:       dto.Stock,
	})
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
