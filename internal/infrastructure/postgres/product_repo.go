package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	productapp "github.com/OmarAraby/go-ecommerce/internal/application/services/product"
	"github.com/OmarAraby/go-ecommerce/internal/domain"
	"github.com/OmarAraby/go-ecommerce/internal/domain/entities"
	"github.com/OmarAraby/go-ecommerce/internal/infrastructure/db"
)

var _ productapp.Repository = (*ProductRepo)(nil)

type ProductRepo struct {
	q    *db.Queries
	pool *pgxpool.Pool
}

func NewProductRepo(pool *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{q: db.New(pool), pool: pool}
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

func (r *ProductRepo) List(ctx context.Context, params productapp.ListParams) ([]*entities.Product, int, error) {
	// --- build WHERE clause dynamically ---
	conditions := []string{}
	args := []any{}
	idx := 1

	if params.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", idx))
		args = append(args, "%"+params.Name+"%")
		idx++
	}
	if params.MinPrice > 0 {
		conditions = append(conditions, fmt.Sprintf("price >= $%d", idx))
		args = append(args, params.MinPrice)
		idx++
	}
	if params.MaxPrice > 0 {
		conditions = append(conditions, fmt.Sprintf("price <= $%d", idx))
		args = append(args, params.MaxPrice)
		idx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	// --- COUNT (total matching rows) ---
	var total int
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM products %s", where)
	if err := r.pool.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("List count: %w", err)
	}

	// --- whitelist sort column and order direction (prevent SQL injection) ---
	sortCol := map[string]string{
		"name": "name", "price": "price", "created_at": "created_at",
	}[params.Sort]
	if sortCol == "" {
		sortCol = "created_at"
	}
	order := "ASC"
	if strings.ToLower(params.Order) == "desc" {
		order = "DESC"
	}

	// --- main query with ORDER BY, LIMIT, OFFSET ---
	offset := (params.Page - 1) * params.Limit
	args = append(args, params.Limit, offset)
	mainSQL := fmt.Sprintf(`
		SELECT id, name, description, price, stock, created_at, updated_at
		FROM products %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d`,
		where, sortCol, order, idx, idx+1,
	)

	rows, err := r.pool.Query(ctx, mainSQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("List query: %w", err)
	}
	defer rows.Close()

	var products []*entities.Product
	for rows.Next() {
		var p entities.Product
		var stock int32
		var createdAt, updatedAt pgtype.Timestamptz
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &stock, &createdAt, &updatedAt); err != nil {
			return nil, 0, fmt.Errorf("List scan: %w", err)
		}
		p.Stock = int(stock)
		p.CreatedAt = createdAt.Time
		p.UpdatedAt = updatedAt.Time
		products = append(products, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("List rows: %w", err)
	}

	return products, total, nil
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
