package product

import (
	"context"
	"fmt"
	"time"

	"github.com/OmarAraby/go-ecommerce/internal/domain"
	"github.com/OmarAraby/go-ecommerce/internal/domain/entities"
)

const maxImagesPerProduct = 6

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
	imgs, err := s.repo.GetImages(ctx, id)
	if err != nil {
		return nil, err
	}
	r := toResponse(p, imgs)
	return &r, nil
}

func (s *service) List(ctx context.Context, params ListParams) ([]ProductResponseDTO, int, error) {
	items, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return []ProductResponseDTO{}, total, nil
	}
	// batch-fetch all images in one query
	ids := make([]int64, len(items))
	for i, p := range items {
		ids[i] = p.ID
	}
	allImgs, err := s.repo.GetImagesByProductIDs(ctx, ids)
	if err != nil {
		return nil, 0, err
	}
	// group images by product_id
	imgMap := make(map[int64][]*entities.ProductImage)
	for _, img := range allImgs {
		imgMap[img.ProductID] = append(imgMap[img.ProductID], img)
	}
	result := make([]ProductResponseDTO, len(items))
	for i, p := range items {
		result[i] = toResponse(p, imgMap[p.ID])
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
	r := toResponse(p, nil)
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
	imgs, err := s.repo.GetImages(ctx, dto.ID)
	if err != nil {
		return nil, err
	}
	r := toResponse(p, imgs)
	return &r, nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) UploadImage(ctx context.Context, dto UploadImageInputDTO) (*ProductResponseDTO, error) {
	// verify product exists
	if _, err := s.repo.GetByID(ctx, dto.ProductID); err != nil {
		return nil, err
	}

	// enforce max 6 images
	count, err := s.repo.CountImages(ctx, dto.ProductID)
	if err != nil {
		return nil, err
	}
	if count >= maxImagesPerProduct {
		return nil, domain.ErrImageLimitReached
	}

	// validate content type
	allowed := map[string]string{
		"image/jpeg": "jpg", "image/png": "png",
		"image/webp": "webp", "image/gif": "gif",
	}
	ext, ok := allowed[dto.ContentType]
	if !ok {
		return nil, domain.ErrInvalidFileType
	}

	// unique filename
	filename := fmt.Sprintf("%d_%d.%s", dto.ProductID, time.Now().UnixNano(), ext)

	// save file
	url, err := s.storage.Save(ctx, "products", filename, dto.File)
	if err != nil {
		return nil, fmt.Errorf("save image: %w", err)
	}

	// first image is automatically main
	isMain := count == 0
	if _, err := s.repo.AddImage(ctx, dto.ProductID, url, isMain); err != nil {
		return nil, err
	}

	return s.productWithImages(ctx, dto.ProductID)
}

func (s *service) DeleteImage(ctx context.Context, productID, imageID int64) (*ProductResponseDTO, error) {
	img, err := s.repo.GetImage(ctx, productID, imageID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.DeleteImage(ctx, productID, imageID); err != nil {
		return nil, err
	}

	// if deleted image was main, promote the first remaining image
	if img.IsMain {
		remaining, _ := s.repo.GetImages(ctx, productID)
		if len(remaining) > 0 {
			_ = s.repo.SetMainImage(ctx, productID, remaining[0].ID)
		}
	}

	return s.productWithImages(ctx, productID)
}

func (s *service) SetMainImage(ctx context.Context, productID, imageID int64) (*ProductResponseDTO, error) {
	// verify image belongs to product
	if _, err := s.repo.GetImage(ctx, productID, imageID); err != nil {
		return nil, err
	}
	if err := s.repo.SetMainImage(ctx, productID, imageID); err != nil {
		return nil, err
	}
	return s.productWithImages(ctx, productID)
}

// productWithImages is a helper that returns the full product DTO with images.
func (s *service) productWithImages(ctx context.Context, productID int64) (*ProductResponseDTO, error) {
	p, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	imgs, err := s.repo.GetImages(ctx, productID)
	if err != nil {
		return nil, err
	}
	r := toResponse(p, imgs)
	return &r, nil
}

func toResponse(p *entities.Product, imgs []*entities.ProductImage) ProductResponseDTO {
	images := make([]ProductImageDTO, len(imgs))
	for i, img := range imgs {
		images[i] = ProductImageDTO{ID: img.ID, URL: img.URL, IsMain: img.IsMain}
	}
	return ProductResponseDTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Images:      images,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
