package service

import (
	"context"
	"errors"
	"testing"

	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type fakeProductRepository struct {
	receivedFilter   model.ProductFilter
	getProductsCalls int
	totalItems       int
}

func (f *fakeProductRepository) GetProducts(ctx context.Context, filter model.ProductFilter) ([]model.Product, error) {
	f.receivedFilter = filter
	f.getProductsCalls++
	return []model.Product{}, nil
}

func (f *fakeProductRepository) CountProducts(ctx context.Context, filter model.ProductFilter) (int, error) {
	return f.totalItems, nil
}

func (f *fakeProductRepository) GetProductByID(ctx context.Context, id int) (model.Product, error) {
	return model.Product{}, nil
}

func (f *fakeProductRepository) CreateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	return model.Product{}, nil
}

func (f *fakeProductRepository) UpdateProduct(
	ctx context.Context,
	product model.Product,
) (model.Product, error) {
	return model.Product{}, nil
}

func (f *fakeProductRepository) UpdateProductStatus(
	ctx context.Context,
	id int,
	status model.ProductStatus,
) (model.Product, error) {
	return model.Product{}, nil
}

func (f *fakeProductRepository) DeleteProduct(
	ctx context.Context,
	id int,
) error {
	return nil
}

func TestProductService_GetProducts_AppliesDefaults(t *testing.T) {
	repository := &fakeProductRepository{}

	service := NewProductService(repository)

	filter := model.ProductFilter{}

	_, err := service.GetProducts(
		context.Background(),
		filter,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repository.receivedFilter.Pagination.Page != 1 {
		t.Errorf(
			"expected page 1, got %d",
			repository.receivedFilter.Pagination.Page,
		)
	}

	if repository.receivedFilter.Pagination.PageSize != 20 {
		t.Errorf(
			"expected page size 20, got %d",
			repository.receivedFilter.Pagination.PageSize,
		)
	}

	if repository.receivedFilter.Status != model.ProductStatusActive {
		t.Errorf(
			"expected status %s, got %s",
			model.ProductStatusActive,
			repository.receivedFilter.Status,
		)
	}

	if repository.receivedFilter.SortBy != "id" {
		t.Errorf(
			"expected sort field id, got %s",
			repository.receivedFilter.SortBy,
		)
	}

	if repository.receivedFilter.SortOrder != "ASC" {
		t.Errorf(
			"expected sort order ASC, got %s",
			repository.receivedFilter.SortOrder,
		)
	}
}

func TestProductService_GetProducts_ReturnsErrorForInvalidStatus(
	t *testing.T,
) {
	repository := &fakeProductRepository{}
	service := NewProductService(repository)

	filter := model.ProductFilter{
		Status: "BANANA",
	}

	_, err := service.GetProducts(
		context.Background(),
		filter,
	)

	if !errors.Is(err, ErrInvalidProductStatus) {
		t.Fatalf(
			"expected ErrInvalidProductStatus, got %v",
			err,
		)
	}

	if repository.getProductsCalls != 0 {
		t.Errorf(
			"expected repository not to be called, got %d calls",
			repository.getProductsCalls,
		)
	}
}

func TestProductService_GetProducts_CalculatesTotalPages(
	t *testing.T,
) {
	repository := &fakeProductRepository{
		totalItems: 21,
	}

	service := NewProductService(repository)

	filter := model.ProductFilter{
		Pagination: model.Pagination{
			Page:     1,
			PageSize: 20,
		},
	}

	productList, err := service.GetProducts(
		context.Background(),
		filter,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if productList.Pagination.TotalItems != 21 {
		t.Errorf(
			"expected 21 total items, got %d",
			productList.Pagination.TotalItems,
		)
	}

	if productList.Pagination.TotalPages != 2 {
		t.Errorf(
			"expected 2 total pages, got %d",
			productList.Pagination.TotalPages,
		)
	}
}

func TestProductService_GetProducts_CapsPageSizeAt100(
	t *testing.T,
) {
	repository := &fakeProductRepository{}
	service := NewProductService(repository)

	filter := model.ProductFilter{
		Pagination: model.Pagination{
			Page:     1,
			PageSize: 500,
		},
	}

	productList, err := service.GetProducts(
		context.Background(),
		filter,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repository.receivedFilter.Pagination.PageSize != 100 {
		t.Errorf(
			"expected repository to receive page size 100, got %d",
			repository.receivedFilter.Pagination.PageSize,
		)
	}

	if productList.Pagination.PageSize != 100 {
		t.Errorf(
			"expected response page size 100, got %d",
			productList.Pagination.PageSize,
		)
	}
}
