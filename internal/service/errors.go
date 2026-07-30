package service

import "errors"

var (
	ErrInvalidProductStatus       = errors.New("invalid product status")
	ErrInvalidSortBy              = errors.New("invalid sort field")
	ErrInvalidSortOrder           = errors.New("invalid sort order")
	ErrInvalidProductQuantity     = errors.New("invalid product quantity")
	ErrInvalidProductMinimumStock = errors.New("invalid product minimum stock")
	ErrInsufficientStock          = errors.New("insufficient stock")
)
