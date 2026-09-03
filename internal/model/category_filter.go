package model

type CategoryFilter struct {
	Pagination Pagination
	Search     string
	SortBy     string
	SortOrder  string
}
