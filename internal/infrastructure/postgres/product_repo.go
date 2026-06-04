package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	productapp "github.com/OmarAraby/go-ecommerce/internal/application/services/product"
	"github.com/OmarAraby/go-ecommerce/internal/domain"
	"github.com/OmarAraby/go-ecommerce/internal/domain/entities"
	"github.com/OmarAraby/go-ecommerce/internal/infrastructure/db"
)

var _ productapp.Repository = (*ProductRepo)(nil)

type ProductRepo struct {
	q *db.Queries
}

func NewProductRepo(pool *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{q: db.New(pool)}
}

func (r *ProductRepo) GetByID(ctx context.Context, id int64) (*entities.Product, error) {
	row, err := r.q.GetProduct(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("GetByID: %w", err)
	}
	return toDomain(row), nil
}

func (r *ProductRepo) List(ctx context.Context) ([]*entities.Product, error) {
	rows, err := r.q.ListProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("List: %w", err)
	}
	products := make([]*entities.Product, len(rows))
	for i, row := range rows {
		products[i] = toDomain(row)
	}
	return products, nil
}

func (r *ProductRepo) Create(ctx context.Context, p *entities.Product) (*entities.Product, error) {
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

func (r *ProductRepo) Update(ctx context.Context, p *entities.Product) (*entities.Product, error) {
	row, err := r.q.UpdateProduct(ctx, db.UpdateProductParams{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       int32(p.Stock),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
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

func toDomain(p db.Product) *entities.Product {
	return &entities.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       int(p.Stock),
		CreatedAt:   p.CreatedAt.Time,
		UpdatedAt:   p.UpdatedAt.Time,
	}
}
