package model

import "time"

type Category struct {
	ID           uint
	Name         string
	ProductCount int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CategoryList struct {
	Categories []Category
	Pagination Pagination
}
