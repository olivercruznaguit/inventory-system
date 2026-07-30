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

func (sr *StockMovementRepository) GetByProductID(ctx context.Context, productID int) ([]model.StockMovement, error) {
	var stockMovements []model.StockMovement

	rows, err := sr.db.Query(ctx, `
	SELECT 
		id,
		type,
		product_id,
		quantity,
		remaining_quantity,
		reason,
		created_at
	FROM
		stock_movements
	WHERE product_id = $1
	ORDER BY created_at DESC
	`, productID)

	if err != nil {
		return nil, fmt.Errorf("get stock movements by product id: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var stockMovement model.StockMovement

		err := rows.Scan(
			&stockMovement.ID,
			&stockMovement.Type,
			&stockMovement.ProductID,
			&stockMovement.Quantity,
			&stockMovement.RemainingQuantity,
			&stockMovement.Reason,
			&stockMovement.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("get stock movements by product id: %w", err)
		}

		stockMovements = append(stockMovements, stockMovement)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get stock movements by product id: %w", err)
	}

	return stockMovements, nil
}
