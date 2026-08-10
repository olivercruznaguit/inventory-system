package testutil

import (
	"context"
	"testing"

	"github.com/olivercruznaguit/inventory-system/internal/database"
	"github.com/olivercruznaguit/inventory-system/internal/model"
)

func SeedCategory(t *testing.T, ctx context.Context, db *database.Database, category model.Category) model.Category {
	t.Helper()

	var insertedCategory model.Category
	err := db.DB().QueryRow(ctx, `
		INSERT INTO categories (name)
		VALUES ($1) RETURNING
		id,
		name,
		created_at,
		updated_at
	`, category.Name).Scan(
		&insertedCategory.ID,
		&insertedCategory.Name,
		&insertedCategory.CreatedAt,
		&insertedCategory.UpdatedAt,
	)

	if err != nil {
		t.Fatalf("failed to seed category: %v", err)
	}

	return insertedCategory
}

func SeedProduct(t *testing.T, ctx context.Context, db *database.Database, product model.Product) model.Product {
	t.Helper()
	var insertedProduct model.Product

	if product.Category == nil {
		t.Fatal("SeedProduct: Category is required")
	}

	categoryID := product.Category.ID

	err := db.DB().QueryRow(ctx, `
        INSERT INTO products (name, price, status, category_id, quantity, minimum_stock)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING 
			id, 
			name, 
			price, 
			status, 
			quantity, 
			minimum_stock, 
			created_at, 
			updated_at
    `, product.Name, product.Price, product.Status, categoryID, product.Quantity, product.MinimumStock).Scan(
		&insertedProduct.ID,
		&insertedProduct.Name,
		&insertedProduct.Price,
		&insertedProduct.Status,
		&insertedProduct.Quantity,
		&insertedProduct.MinimumStock,
		&insertedProduct.CreatedAt,
		&insertedProduct.UpdatedAt,
	)

	insertedProduct.Category = product.Category

	if err != nil {
		t.Fatalf("failed to seed product: %v", err)
	}

	return insertedProduct
}

func SeedUser(t *testing.T, ctx context.Context, db *database.Database, user model.User) model.User {
	t.Helper()
	var insertedUser model.User

	err := db.DB().QueryRow(ctx, `
        INSERT INTO users (email, password_hash)
        VALUES ($1, $2)
        RETURNING 
			id, 
			email, 
			password_hash,
			created_at, 
			updated_at
    `, user.Email, user.PasswordHash).Scan(
		&insertedUser.ID,
		&insertedUser.Email,
		&insertedUser.PasswordHash,
		&insertedUser.CreatedAt,
		&insertedUser.UpdatedAt,
	)

	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	return insertedUser
}
