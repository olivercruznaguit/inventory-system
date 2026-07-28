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
