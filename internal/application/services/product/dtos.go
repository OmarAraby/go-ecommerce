package product

import (
	"io"
	"time"
)

type CreateProductDTO struct {
	Name        string
	Description string
	Price       float64
	Stock       int
}

type UpdateProductDTO struct {
	ID          int64
	Name        string
	Description string
	Price       float64
	Stock       int
}

// ListParams holds pagination, sorting, and filter options for product listing.
type ListParams struct {
	Page     int
	Limit    int
	Sort     string  // "name" | "price" | "created_at"
	Order    string  // "asc" | "desc"
	Name     string  // partial match (ILIKE)
	MinPrice float64
	MaxPrice float64
}

// ProductResponseDTO is what the API layer sees — never the raw domain entity.
type ProductResponseDTO struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UploadImageInputDTO struct {
	ProductID   int64
	File        io.Reader
	Filename    string // original filename (used for extension)
	ContentType string // pre-detected by handler
}
