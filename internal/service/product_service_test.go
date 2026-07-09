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
}

func (f *fakeProductRepository) GetProducts(ctx context.Context, filter model.ProductFilter) ([]model.Product, error) {
	f.receivedFilter = filter
	f.getProductsCalls++
	return []model.Product{}, nil
}

func (f *fakeProductRepository) CountProducts(ctx context.Context, filter model.ProductFilter) (int, error) {
	return 0, nil
}

func (f *fakeProductRepository) GetProductByID(ctx context.Context, id int) (model.Product, error) {
	return model.Product{}, nil
}

func (f *fakeProductRepository) CreateProduct(
	ctx context.Context,
	product model.Product,
) (model.Product, error) {
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
