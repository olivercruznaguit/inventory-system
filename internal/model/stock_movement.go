package model

import "time"

type StockMovementType string

const (
	StockMovementTypeIn  StockMovementType = "IN"
	StockMovementTypeOut StockMovementType = "OUT"
)

type StockMovement struct {
	ID                int64
	ProductID         int64
	Type              StockMovementType
	Quantity          int
	RemainingQuantity int
	Reason            *string
	CreatedAt         time.Time
}
