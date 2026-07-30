package repository

import (
	"context"
	"fmt"

	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type StockMovementRepository struct {
	db DBTX
}

func NewStockMovementRepository(db DBTX) *StockMovementRepository {
	return &StockMovementRepository{
		db: db,
	}
}

func (sr *StockMovementRepository) Create(ctx context.Context, movement *model.StockMovement) error {
	_, err := sr.db.Exec(ctx, `
	INSERT INTO stock_movements (
		product_id,
		type,
		quantity,
		remaining_quantity,
		reason
	)
	VALUES ($1, $2, $3, $4, $5)
    `,
		movement.ProductID,
		movement.Type,
		movement.Quantity,
		movement.RemainingQuantity,
		movement.Reason)

	if err != nil {
		return fmt.Errorf("create stock movement: %w", err)
	}

	return nil
}
