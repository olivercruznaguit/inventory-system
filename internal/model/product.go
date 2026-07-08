package model

type Product struct {
	ID     uint
	Name   string
	Price  float64
	Status ProductStatus
}

type ProductList struct {
	Products   []Product
	Pagination Pagination
}
