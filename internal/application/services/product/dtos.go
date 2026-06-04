package product

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
