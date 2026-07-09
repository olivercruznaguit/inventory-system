package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/olivercruznaguit/inventory-system/internal/model"
	"github.com/olivercruznaguit/inventory-system/internal/repository"
)

type ProductService struct {
	repository *repository.ProductRepository
}

func NewProductService(repository *repository.ProductRepository) *ProductService {
	return &ProductService{
		repository: repository,
	}
}

func (ps *ProductService) GetProducts(ctx context.Context, filter model.ProductFilter) (model.ProductList, error) {
	if filter.Pagination.Page < 1 {
		filter.Pagination.Page = 1
	}

	if filter.Pagination.PageSize < 1 {
		filter.Pagination.PageSize = 20
	}

	if filter.Pagination.PageSize > 100 {
		filter.Pagination.PageSize = 100
	}

	if filter.SortBy == "" {
		filter.SortBy = "id"
	}

	switch filter.SortBy {
	case "id", "name", "price":
		filter.SortBy = strings.ToLower(filter.SortBy)
	default:
		return model.ProductList{},
			fmt.Errorf("get products: %w", ErrInvalidSortBy)
	}

	if filter.SortOrder == "" {
		filter.SortOrder = "asc"
	}

	switch strings.ToLower(filter.SortOrder) {
	case "asc", "desc":
		filter.SortOrder = strings.ToUpper(filter.SortOrder)
	default:
		return model.ProductList{},
			fmt.Errorf("get products: %w", ErrInvalidSortOrder)
	}

	if filter.Status == "" {
		filter.Status = model.ProductStatusActive
	}

	if !filter.Status.IsValid() {
		return model.ProductList{}, fmt.Errorf("get products: %w", ErrInvalidProductStatus)
	}

	products, err := ps.repository.GetProducts(ctx, filter)
	if err != nil {
		return model.ProductList{}, fmt.Errorf("get products: %w", err)
	}

	totalItems, err := ps.repository.CountProducts(ctx, filter)
	if err != nil {
		return model.ProductList{}, fmt.Errorf("count products: %w", err)
	}

	totalPages := (totalItems + filter.Pagination.PageSize - 1) / filter.Pagination.PageSize

	filter.Pagination.TotalItems = totalItems
	filter.Pagination.TotalPages = totalPages

	return model.ProductList{
		Products:   products,
		Pagination: filter.Pagination,
	}, nil
}

func (ps *ProductService) GetProductByID(ctx context.Context, id int) (model.Product, error) {
	return ps.repository.GetProductByID(ctx, id)
}

func (ps *ProductService) CreateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	return ps.repository.CreateProduct(ctx, product)
}

func (ps *ProductService) UpdateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	if !product.Status.IsValid() {
		return model.Product{}, fmt.Errorf("updated product: %w", ErrInvalidProductStatus)
	}

	updatedProduct, err := ps.repository.UpdateProduct(ctx, product)
	if err != nil {
		return model.Product{}, fmt.Errorf("update product: %w", err)
	}
	return updatedProduct, nil
}

func (ps *ProductService) UpdateProductStatus(ctx context.Context, id int, status model.ProductStatus) (model.Product, error) {
	if !status.IsValid() {
		return model.Product{}, fmt.Errorf("update product status: %w", ErrInvalidProductStatus)
	}

	updatedProduct, err := ps.repository.UpdateProductStatus(ctx, id, status)
	if err != nil {
		return model.Product{}, fmt.Errorf("update product status: %w", err)
	}

	return updatedProduct, nil

}

func (ps *ProductService) DeleteProduct(ctx context.Context, id int) error {
	return ps.repository.DeleteProduct(ctx, id)
}
