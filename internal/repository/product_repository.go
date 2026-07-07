package repository

import (
	"context"

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
