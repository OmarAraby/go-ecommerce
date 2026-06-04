package product

import (
	"context"
	"fmt"
	"time"

	"github.com/OmarAraby/go-ecommerce/internal/domain"
	"github.com/OmarAraby/go-ecommerce/internal/domain/entities"
)

var _ Service = (*service)(nil)

type service struct {
	repo    Repository
	storage FileStorage
}

func NewService(repo Repository, storage FileStorage) Service {
	return &service{repo: repo, storage: storage}
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

func (s *service) UploadImage(ctx context.Context, dto UploadImageInputDTO) (*ProductResponseDTO, error) {
	// 1. verify product exists
	if _, err := s.repo.GetByID(ctx, dto.ProductID); err != nil {
		return nil, err
	}

	// 2. validate content type
	allowed := map[string]string{
		"image/jpeg": "jpg",
		"image/png":  "png",
		"image/webp": "webp",
		"image/gif":  "gif",
	}
	ext, ok := allowed[dto.ContentType]
	if !ok {
		return nil, domain.ErrInvalidFileType
	}

	// 3. unique filename: {productID}_{nanoseconds}.{ext}
	filename := fmt.Sprintf("%d_%d.%s", dto.ProductID, time.Now().UnixNano(), ext)

	// 4. persist file via storage
	url, err := s.storage.Save(ctx, "products", filename, dto.File)
	if err != nil {
		return nil, fmt.Errorf("save image: %w", err)
	}

	// 5. update product record
	p, err := s.repo.UpdateImage(ctx, dto.ProductID, url)
	if err != nil {
		return nil, err
	}
	r := toResponse(p)
	return &r, nil
}

func toResponse(p *entities.Product) ProductResponseDTO {
	return ProductResponseDTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		ImageURL:    p.ImageURL,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
