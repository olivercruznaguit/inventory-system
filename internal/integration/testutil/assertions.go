package testutil

import (
	"context"
	"testing"

	"github.com/olivercruznaguit/inventory-system/internal/database"
	"github.com/olivercruznaguit/inventory-system/internal/model"
)

func AssertProductQuantity(t *testing.T, ctx context.Context, db *database.Database, productID uint, expected int) {
	t.Helper()

	var quantity int

	row := db.DB().QueryRow(ctx, `SELECT quantity FROM products WHERE id = $1`, productID)

	err := row.Scan(&quantity)
	if err != nil {
		t.Fatalf("failed to scan product quantity: %v", err)
	}

	if quantity != expected {
		t.Fatalf("expected quantity %d, got %d", expected, quantity)
	}

}

func AssertStockMovement(t *testing.T, ctx context.Context, db *database.Database, productID uint, expected model.StockMovement) {
	t.Helper()

	var quantity int
	var movementType model.StockMovementType
	var remainingQuantity int

	row := db.DB().QueryRow(ctx, `
	SELECT
		type,
		quantity,
		remaining_quantity
	FROM 
		stock_movements
	WHERE 
		product_id = $1 
	ORDER BY 
		created_at DESC 
	LIMIT 1`, productID)

	err := row.Scan(&movementType, &quantity, &remainingQuantity)

	if err != nil {
		t.Fatalf("failed to query stock movements: %v", err)
	}

	if movementType != expected.Type {
		t.Fatalf("expected movement type %v, got %v", expected.Type, movementType)
	}

	if quantity != expected.Quantity {
		t.Fatalf("expected movement quantity %d, got %d", expected.Quantity, quantity)
	}

	if remainingQuantity != expected.RemainingQuantity {
		t.Fatalf("expected remaining quantity %d, got %d", expected.RemainingQuantity, remainingQuantity)
	}

}

func AssertStockMovementCount(t *testing.T, ctx context.Context, db *database.Database, productID uint, expected int) {
	t.Helper()

	var count int

	row := db.DB().QueryRow(ctx, `SELECT COUNT(*) FROM stock_movements WHERE product_id = $1`, productID)

	err := row.Scan(&count)
	if err != nil {
		t.Fatalf("failed to scan stock movement count: %v", err)
	}

	if count != expected {
		t.Fatalf("expected stock movement count %d, got %d", expected, count)
	}
}
