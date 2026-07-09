package model

import "time"

type Product struct {
	ID        uint
	Name      string
	Price     float64
	Status    ProductStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProductList struct {
	Products   []Product
	Pagination Pagination
}
