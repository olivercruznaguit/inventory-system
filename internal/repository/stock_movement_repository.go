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

func (sr *StockMovementRepository) GetRecent(ctx context.Context, limit int) ([]model.RecentStockMovement, error) {
	var recentStockMovements []model.RecentStockMovement

	rows, err := sr.db.Query(ctx,
		`SELECT
			sm.id,
			sm.product_id,
			p.name,
			sm.type,
			sm.quantity,
			sm.remaining_quantity,
			sm.reason,
			sm.created_at
		FROM 
			stock_movements sm
		JOIN 
			products p ON p.id = sm.product_id
		ORDER BY 
			sm.created_at DESC
		LIMIT $1
		`,
		limit)

	if err != nil {
		return nil, fmt.Errorf("get recent stock movements: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var recentStockMovement model.RecentStockMovement

		err := rows.Scan(
			&recentStockMovement.ID,
			&recentStockMovement.ProductID,
			&recentStockMovement.ProductName,
			&recentStockMovement.Type,
			&recentStockMovement.Quantity,
			&recentStockMovement.RemainingQuantity,
			&recentStockMovement.Reason,
			&recentStockMovement.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("get recent stock movements: %w", err)
		}

		recentStockMovements = append(recentStockMovements, recentStockMovement)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get recent stock movements: %w", err)
	}

	return recentStockMovements, nil
}
