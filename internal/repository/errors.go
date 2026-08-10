package repository

import "errors"

var (
	ErrProductNotFound       = errors.New("product not found")
	ErrCategoryAlreadyExists = errors.New("category already exist")
	ErrCategoryNotFound      = errors.New("category not found")
	ErrUserNotFound          = errors.New("user not found")
)
