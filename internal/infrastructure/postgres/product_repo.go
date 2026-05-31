package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/OmarAraby/go-ecommerce/internal/application/interfaces/repositories"
	"github.com/OmarAraby/go-ecommerce/internal/db"
	"github.com/OmarAraby/go-ecommerce/internal/domain"
)

var _ repositories.ProductRepository = (*ProductRepo)(nil)

// ProductRepo implements application.ProductRepository using sqlc-generated queries.
type ProductRepo struct {
	q *db.Queries
}

func NewProductRepo(pool *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{q: db.New(pool)}
}

func (r *ProductRepo) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	row, err := r.q.GetProduct(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetByID: %w", err)
	}
	return toDomain(row), nil
}

func (r *ProductRepo) List(ctx context.Context) ([]*domain.Product, error) {
	rows, err := r.q.ListProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("List: %w", err)
	}
	products := make([]*domain.Product, len(rows))
	for i, row := range rows {
		products[i] = toDomain(row)
	}
	return products, nil
}

func (r *ProductRepo) Create(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	row, err := r.q.CreateProduct(ctx, db.CreateProductParams{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       int32(p.Stock),
	})
	if err != nil {
		return nil, fmt.Errorf("Create: %w", err)
	}
	return toDomain(row), nil
}

func (r *ProductRepo) Update(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	row, err := r.q.UpdateProduct(ctx, db.UpdateProductParams{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       int32(p.Stock),
	})
	if err != nil {
		return nil, fmt.Errorf("Update: %w", err)
	}
	return toDomain(row), nil
}

func (r *ProductRepo) Delete(ctx context.Context, id int64) error {
	if err := r.q.DeleteProduct(ctx, id); err != nil {
		return fmt.Errorf("Delete: %w", err)
	}
	return nil
}

// toDomain maps the sqlc-generated type to the domain entity.
// The infrastructure layer owns this mapping — the domain stays pure.
func toDomain(p db.Product) *domain.Product {
	return &domain.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       int(p.Stock),
		CreatedAt:   p.CreatedAt.Time,
		UpdatedAt:   p.UpdatedAt.Time,
	}
}
