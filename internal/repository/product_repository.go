package repository

import (
	"context"
	"errors"
	"fmt"

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

func (pr *ProductRepository) GetProducts(ctx context.Context, pagination model.Pagination) ([]model.Product, error) {

	offset := (pagination.Page - 1) * pagination.PageSize
	limit := pagination.PageSize
	rows, err := pr.db.Query(ctx, `
        SELECT
            id,
            name,
            price
        FROM products
        ORDER BY id
		LIMIT $1 OFFSET $2
    `, limit, offset)
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
            price
        FROM products
        WHERE id = $1
    `, id)

	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Price,
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
        RETURNING id, name, price
    `, product.Name, product.Price).Scan(
		&createdProduct.ID,
		&createdProduct.Name,
		&createdProduct.Price,
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
        SET name = $1, price = $2
        WHERE id = $3
        RETURNING id, name, price
    `, product.Name, product.Price, product.ID).Scan(
		&updatedProduct.ID,
		&updatedProduct.Name,
		&updatedProduct.Price,
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
