package model

import "time"

type Product struct {
	ID           uint
	Name         string
	Price        float64
	Status       ProductStatus
	Quantity     int
	MinimumStock int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Category     *Category
}

type ProductList struct {
	Products   []Product
	Pagination Pagination
}

type RecentStockMovement struct {
	ID                int64
	ProductID         int64
	ProductName       string
	Type              StockMovementType
	Quantity          int
	RemainingQuantity int
	Reason            *string
	CreatedAt         time.Time
}
