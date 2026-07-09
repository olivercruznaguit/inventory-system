package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (pr *ProductRepository) GetProducts(ctx context.Context, filter model.ProductFilter) ([]model.Product, error) {
	pagination := filter.Pagination
	offset := (pagination.Page - 1) * pagination.PageSize
	limit := pagination.PageSize

	var queryParts []string
	var conditions []string
	var args []any

	queryParts = append(queryParts, "SELECT id, name, price, status FROM products")

	if filter.Status != "" {
		args = append(args, filter.Status)

		conditions = append(
			conditions,
			fmt.Sprintf("status = $%d", len(args)),
		)
	}

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")

		conditions = append(
			conditions,
			fmt.Sprintf("name ILIKE $%d", len(args)),
		)
	}

	if len(conditions) > 0 {
		queryParts = append(
			queryParts,
			"WHERE "+strings.Join(conditions, " AND "),
		)
	}

	queryParts = append(queryParts, "ORDER BY id")

	args = append(args, limit)
	queryParts = append(queryParts, fmt.Sprintf("LIMIT $%d", len(args)))

	args = append(args, offset)
	queryParts = append(queryParts, fmt.Sprintf("OFFSET $%d", len(args)))

	queryString := strings.Join(queryParts, " ")

	rows, err := pr.db.Query(ctx, queryString, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []model.Product

	for rows.Next() {
		var product model.Product

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("get products: %w", err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get products: %w", err)
	}

	return products, nil
}

func (pr *ProductRepository) GetProductByID(ctx context.Context, id int) (model.Product, error) {
	var product model.Product

	row := pr.db.QueryRow(ctx, `
        SELECT
            id,
            name,
            price,
			status
        FROM products
        WHERE id = $1
    `, id)

	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Status,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Product{}, ErrProductNotFound
		}
		return model.Product{}, fmt.Errorf("get product: %w", err)
	}

	return product, nil
}

func (pr *ProductRepository) CreateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	var createdProduct model.Product

	err := pr.db.QueryRow(ctx, `
        INSERT INTO products (name, price)
        VALUES ($1, $2)
        RETURNING id, name, price, status
    `, product.Name, product.Price).Scan(
		&createdProduct.ID,
		&createdProduct.Name,
		&createdProduct.Price,
		&createdProduct.Status,
	)
	if err != nil {
		return model.Product{}, fmt.Errorf("create product: %w", err)
	}

	return createdProduct, nil
}

func (pr *ProductRepository) UpdateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	var updatedProduct model.Product

	err := pr.db.QueryRow(ctx, `
        UPDATE products
        SET name = $1, price = $2, status = $3
        WHERE id = $4
        RETURNING id, name, price, status
    `, product.Name, product.Price, product.Status, product.ID).Scan(
		&updatedProduct.ID,
		&updatedProduct.Name,
		&updatedProduct.Price,
		&updatedProduct.Status,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Product{}, ErrProductNotFound
		}
		return model.Product{}, fmt.Errorf("update product: %w", err)
	}

	return updatedProduct, nil
}

func (pr *ProductRepository) UpdateProductStatus(ctx context.Context, id int, status model.ProductStatus) (model.Product, error) {
	var updatedProduct model.Product

	err := pr.db.QueryRow(ctx, `
        UPDATE products
        SET status = $1
        WHERE id = $2
        RETURNING id, name, price, status
    `, status, id).Scan(
		&updatedProduct.ID,
		&updatedProduct.Name,
		&updatedProduct.Price,
		&updatedProduct.Status,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Product{}, ErrProductNotFound
		}
		return model.Product{}, fmt.Errorf("update product: %w", err)
	}

	return updatedProduct, nil
}

func (pr *ProductRepository) DeleteProduct(ctx context.Context, id int) error {
	result, err := pr.db.Exec(ctx, `
		DELETE FROM products
		WHERE id = $1
	`, id)

	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrProductNotFound
	}

	return nil
}

func (pr *ProductRepository) CountProducts(ctx context.Context, filter model.ProductFilter) (int, error) {
	var count int
	var queryParts []string
	var conditions []string
	var args []any

	queryParts = append(queryParts, "SELECT COUNT(*) FROM products")

	if filter.Status != "" {
		args = append(args, filter.Status)
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)))
	}

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", len(args)))
	}

	if len(conditions) > 0 {
		queryParts = append(
			queryParts,
			"WHERE "+strings.Join(conditions, " AND "),
		)
	}

	queryString := strings.Join(queryParts, " ")

	err := pr.db.QueryRow(ctx, queryString, args...).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("count products: %w", err)
	}

	return count, nil
}
