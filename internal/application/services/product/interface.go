package product

import "context"

type Service interface {
	GetByID(ctx context.Context, id int64) (*ProductResponseDTO, error)
	List(ctx context.Context, params ListParams) ([]ProductResponseDTO, int, error)
	Create(ctx context.Context, dto CreateProductDTO) (*ProductResponseDTO, error)
	Update(ctx context.Context, dto UpdateProductDTO) (*ProductResponseDTO, error)
	Delete(ctx context.Context, id int64) error
	UploadImage(ctx context.Context, dto UploadImageInputDTO) (*ProductResponseDTO, error)
}
