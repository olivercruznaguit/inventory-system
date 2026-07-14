package model

type ProductFilter struct {
	Pagination Pagination
	Search     string
	Status     ProductStatus
	SortBy     string
	SortOrder  string
	CategoryID *uint
}
