package testutil

import (
	"context"
	"testing"

	"github.com/olivercruznaguit/inventory-system/internal/database"
)

func CleanupDatabase(t *testing.T, ctx context.Context, db *database.Database) {
	t.Helper()

	_, err := db.DB().Exec(ctx, `
		TRUNCATE TABLE
			stock_movements,
			products,
			categories
		RESTART IDENTITY CASCADE;`)
	if err != nil {
		t.Fatalf("failed to cleanup database: %v", err)
	}
}
