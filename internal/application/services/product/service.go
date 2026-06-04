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

func (s *service) GetByID(ctx context.Context, id int64) (*ProductResponseDTO, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	r := toResponse(p)
	return &r, nil
}

func (s *service) List(ctx context.Context, params ListParams) ([]ProductResponseDTO, int, error) {
	items, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	result := make([]ProductResponseDTO, len(items))
	for i, item := range items {
		result[i] = toResponse(item)
	}
	return result, total, nil
}

func (s *service) Create(ctx context.Context, dto CreateProductDTO) (*ProductResponseDTO, error) {
	p, err := s.repo.Create(ctx, &entities.Product{
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		Stock:       dto.Stock,
	})
	if err != nil {
		return nil, err
	}
	r := toResponse(p)
	return &r, nil
}

func (s *service) Update(ctx context.Context, dto UpdateProductDTO) (*ProductResponseDTO, error) {
	p, err := s.repo.Update(ctx, &entities.Product{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		Stock:       dto.Stock,
	})
	if err != nil {
		return nil, err
	}
	r := toResponse(p)
	return &r, nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func toResponse(p *entities.Product) ProductResponseDTO {
	return ProductResponseDTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
