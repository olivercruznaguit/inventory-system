package repository

import (
	"context"

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

func (pr *ProductRepository) GetProducts(ctx context.Context) ([]model.Product, error) {
	rows, err := pr.db.Query(ctx, `
        SELECT
            id,
            name,
            price
        FROM products
        ORDER BY id
    `)
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
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
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
		return model.Product{}, err
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
		return model.Product{}, err
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
		return model.Product{}, err
	}

	return updatedProduct, nil
}

func (pr *ProductRepository) DeleteProduct(ctx context.Context, id int) error {
	result, err := pr.db.Exec(ctx, `
		DELETE FROM products
		WHERE id = $1
	`, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
