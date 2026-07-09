package service

import "errors"

var (
	ErrInvalidProductStatus = errors.New("invalid product status")
	ErrInvalidSortBy        = errors.New("invalid sort field")
	ErrInvalidSortOrder     = errors.New("invalid sort order")
)
