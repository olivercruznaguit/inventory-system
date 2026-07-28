package repository

import "errors"

var (
	ErrProductNotFound       = errors.New("product not found")
	ErrCategoryAlreadyExists = errors.New("category already exist")
	ErrCategoryNotFound      = errors.New("category not found")
	ErrInsufficientStock     = errors.New("insufficient stock")
)
