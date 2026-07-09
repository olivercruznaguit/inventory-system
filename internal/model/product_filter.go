package model

type ProductFilter struct {
	Pagination Pagination
	Search     string
	Status     ProductStatus
}
