package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type ProductRepository struct {
	db DBTX
}

func NewProductRepository(db DBTX) *ProductRepository {
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

	queryParts = append(queryParts,
		`SELECT 
			p.id, 
			p.name, 
			p.price, 
			p.status, 
			p.quantity, 
			p.minimum_stock, 
			p.created_at, 
			p.updated_at, 
			c.id, 
			c.name 
		FROM products p LEFT JOIN categories c ON p.category_id = c.id`)

	if filter.Status != "" {
		args = append(args, filter.Status)

		conditions = append(
			conditions,
			fmt.Sprintf("p.status = $%d", len(args)),
		)
	}

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")

		conditions = append(
			conditions,
			fmt.Sprintf("p.name ILIKE $%d", len(args)),
		)
	}

	if filter.CategoryID != nil {
		args = append(args, *filter.CategoryID)

		conditions = append(conditions, fmt.Sprintf("p.category_id = $%d", len(args)))
	}

	if len(conditions) > 0 {
		queryParts = append(
			queryParts,
			"WHERE "+strings.Join(conditions, " AND "),
		)
	}

	queryParts = append(queryParts, fmt.Sprintf("ORDER BY p.%s %s", filter.SortBy, filter.SortOrder))

	args = append(args, limit)
	queryParts = append(queryParts, fmt.Sprintf("LIMIT $%d", len(args)))

	args = append(args, offset)
	queryParts = append(queryParts, fmt.Sprintf("OFFSET $%d", len(args)))

	queryString := strings.Join(queryParts, " ")

	rows, err := pr.db.Query(ctx, queryString, args...)

	if err != nil {
		return nil, fmt.Errorf("get products: %w", err)
	}

	defer rows.Close()

	var products []model.Product

	for rows.Next() {
		var product model.Product
		var categoryID pgtype.Int4
		var categoryName pgtype.Text

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Status,
			&product.Quantity,
			&product.MinimumStock,
			&product.CreatedAt,
			&product.UpdatedAt,
			&categoryID,
			&categoryName,
		)

		if err != nil {
			return nil, fmt.Errorf("get products: %w", err)
		}

		if categoryID.Valid {
			product.Category = &model.Category{
				ID:   uint(categoryID.Int32),
				Name: categoryName.String,
			}

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
	var categoryID pgtype.Int4
	var categoryName pgtype.Text

	row := pr.db.QueryRow(ctx, `
        SELECT
            p.id,
            p.name,
            p.price,
			p.status,
			p.quantity, 
			p.minimum_stock,
			p.created_at,
			p.updated_at,
			c.id,
			c.name 
        FROM products p
		LEFT JOIN categories c
			ON p.category_id = c.id
        WHERE p.id = $1
    `, id)

	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Status,
		&product.Quantity,
		&product.MinimumStock,
		&product.CreatedAt,
		&product.UpdatedAt,
		&categoryID,
		&categoryName,
	)

	if categoryID.Valid {
		product.Category = &model.Category{
			ID:   uint(categoryID.Int32),
			Name: categoryName.String,
		}

	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Product{}, ErrProductNotFound
		}
		return model.Product{}, fmt.Errorf("get product: %w", err)
	}

	return product, nil
}

func (pr *ProductRepository) GetProductByIDForUpdate(ctx context.Context, id int) (model.Product, error) {
	var product model.Product

	row := pr.db.QueryRow(ctx, `
        SELECT
            id,
            name,
            price,
			status,
			quantity, 
			minimum_stock,
			created_at,
			updated_at
        FROM products
		WHERE id = $1 
		FOR UPDATE
    `, id)

	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Status,
		&product.Quantity,
		&product.MinimumStock,
		&product.CreatedAt,
		&product.UpdatedAt,
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
	var categoryID any

	if product.Category != nil {
		categoryID = product.Category.ID
	}

	err := pr.db.QueryRow(ctx, `
        INSERT INTO products (name, price, category_id, minimum_stock)
        VALUES ($1, $2, $3, $4)
        RETURNING 
			id, 
			name, 
			price, 
			status, 
			quantity, 
			minimum_stock, 
			created_at, 
			updated_at
    `, product.Name, product.Price, categoryID, product.MinimumStock).Scan(
		&createdProduct.ID,
		&createdProduct.Name,
		&createdProduct.Price,
		&createdProduct.Status,
		&createdProduct.Quantity,
		&createdProduct.MinimumStock,
		&createdProduct.CreatedAt,
		&createdProduct.UpdatedAt,
	)
	if err != nil {
		return model.Product{}, fmt.Errorf("create product: %w", err)
	}

	return createdProduct, nil
}

func (pr *ProductRepository) UpdateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	var updatedProduct model.Product
	var categoryID any

	if product.Category != nil {
		categoryID = product.Category.ID
	}

	err := pr.db.QueryRow(ctx, `
        UPDATE products
        SET 
			name = $1, 
			price = $2, 
			status = $3, 
			category_id = $4, 
			minimum_stock = $5, 
			updated_at = NOW()
        WHERE id = $6
        RETURNING 
			id, 
			name, 
			price, 
			status, 
			quantity, 
			minimum_stock, 
			created_at, 
			updated_at
    `, product.Name, product.Price, product.Status, categoryID, product.MinimumStock, product.ID).Scan(
		&updatedProduct.ID,
		&updatedProduct.Name,
		&updatedProduct.Price,
		&updatedProduct.Status,
		&updatedProduct.Quantity,
		&updatedProduct.MinimumStock,
		&updatedProduct.CreatedAt,
		&updatedProduct.UpdatedAt,
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
        SET status = $1, updated_at = NOW()
        WHERE id = $2
        RETURNING id, name, price, status, quantity, minimum_stock, created_at, updated_at
    `, status, id).Scan(
		&updatedProduct.ID,
		&updatedProduct.Name,
		&updatedProduct.Price,
		&updatedProduct.Status,
		&updatedProduct.Quantity,
		&updatedProduct.MinimumStock,
		&updatedProduct.CreatedAt,
		&updatedProduct.UpdatedAt,
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

	if filter.CategoryID != nil {
		args = append(args, *filter.CategoryID)

		conditions = append(conditions, fmt.Sprintf("category_id = $%d", len(args)))
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

// func (pr *ProductRepository) StockIn(ctx context.Context, productID int, request model.StockRequest) (model.Product, error) {
// 	var updatedProduct model.Product

// 	err := pr.db.QueryRow(ctx, `
//         UPDATE products
// 		SET quantity = quantity + $1,
// 			updated_at = NOW()
// 		WHERE id = $2
// 		RETURNING
// 			id,
// 			name,
// 			price,
// 			status,
// 			quantity,
// 			minimum_stock,
// 			created_at,
// 			updated_at
// 		`, request.Quantity, productID).Scan(
// 		&updatedProduct.ID,
// 		&updatedProduct.Name,
// 		&updatedProduct.Price,
// 		&updatedProduct.Status,
// 		&updatedProduct.Quantity,
// 		&updatedProduct.MinimumStock,
// 		&updatedProduct.CreatedAt,
// 		&updatedProduct.UpdatedAt,
// 	)

// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			return model.Product{}, ErrProductNotFound
// 		}
// 		return model.Product{}, fmt.Errorf("stock in: %w", err)
// 	}

// 	return updatedProduct, nil
// }

// func (pr *ProductRepository) StockOut(ctx context.Context, productID int, request model.StockRequest) (model.Product, error) {
// 	var updatedProduct model.Product

// 	err := pr.db.QueryRow(ctx, `
//         UPDATE products
// 		SET quantity = quantity - $1,
// 			updated_at = NOW()
// 		WHERE id = $2
// 			AND quantity >= $1
// 		RETURNING
// 			id,
// 			name,
// 			price,
// 			status,
// 			quantity,
// 			minimum_stock,
// 			created_at,
// 			updated_at
// 		`, request.Quantity, productID).Scan(
// 		&updatedProduct.ID,
// 		&updatedProduct.Name,
// 		&updatedProduct.Price,
// 		&updatedProduct.Status,
// 		&updatedProduct.Quantity,
// 		&updatedProduct.MinimumStock,
// 		&updatedProduct.CreatedAt,
// 		&updatedProduct.UpdatedAt,
// 	)

// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			_, err := pr.GetProductByID(ctx, productID)

// 			if err != nil {
// 				if errors.Is(err, ErrProductNotFound) {
// 					return model.Product{}, ErrProductNotFound
// 				}

// 				return model.Product{}, fmt.Errorf("verify product after stock out: %w", err)
// 			}

// 			return model.Product{}, ErrInsufficientStock
// 		}

// 		return model.Product{}, fmt.Errorf("stock out: %w", err)
// 	}

// 	return updatedProduct, nil
// }

func (pr *ProductRepository) UpdateProductQuantity(ctx context.Context, productID int, newQuantity int) error {
	cmdTag, err := pr.db.Exec(ctx, `
       UPDATE products
		SET quantity = $1,
			updated_at = NOW()
		WHERE id = $2;
    `, newQuantity, productID)

	if err != nil {
		return fmt.Errorf("update product quantity: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrProductNotFound
	}

	return nil
}
