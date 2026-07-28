package service

import (
	"context"
	"errors"
	"testing"

	"github.com/olivercruznaguit/inventory-system/internal/model"
	"github.com/olivercruznaguit/inventory-system/internal/repository"
)

type fakeProductRepository struct {
	receivedFilter     model.ProductFilter
	getProductsCalls   int
	createProductCalls int
	totalItems         int
	createdProduct     model.Product
}

type fakeCategoryRepository struct {
	getCategoryByIDCalls int
	category             model.Category
	err                  error
}

func (f *fakeCategoryRepository) GetCategoryByID(ctx context.Context, id int) (model.Category, error) {
	f.getCategoryByIDCalls++

	if f.err != nil {
		return model.Category{}, f.err
	}

	return f.category, nil
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
	f.createProductCalls++
	return f.createdProduct, nil
}

func (f *fakeProductRepository) UpdateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	return model.Product{}, nil
}

func (f *fakeProductRepository) UpdateProductStatus(ctx context.Context, id int, status model.ProductStatus) (model.Product, error) {
	return model.Product{}, nil
}

func (f *fakeProductRepository) DeleteProduct(ctx context.Context, id int) error {
	return nil
}

func (f *fakeProductRepository) StockIn(ctx context.Context, productID int, request model.StockRequest) (model.Product, error) {
	return model.Product{}, nil
}

func (f *fakeProductRepository) StockOut(ctx context.Context, productID int, request model.StockRequest) (model.Product, error) {
	return model.Product{}, nil
}

func TestProductService_GetProducts_AppliesDefaults(t *testing.T) {
	productRepository := &fakeProductRepository{}
	categoryRepository := &fakeCategoryRepository{}

	service := NewProductService(productRepository, categoryRepository)

	filter := model.ProductFilter{}

	_, err := service.GetProducts(
		context.Background(),
		filter,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if productRepository.receivedFilter.Pagination.Page != 1 {
		t.Errorf(
			"expected page 1, got %d",
			productRepository.receivedFilter.Pagination.Page,
		)
	}

	if productRepository.receivedFilter.Pagination.PageSize != 20 {
		t.Errorf(
			"expected page size 20, got %d",
			productRepository.receivedFilter.Pagination.PageSize,
		)
	}

	if productRepository.receivedFilter.Status != model.ProductStatusActive {
		t.Errorf(
			"expected status %s, got %s",
			model.ProductStatusActive,
			productRepository.receivedFilter.Status,
		)
	}

	if productRepository.receivedFilter.SortBy != "id" {
		t.Errorf(
			"expected sort field id, got %s",
			productRepository.receivedFilter.SortBy,
		)
	}

	if productRepository.receivedFilter.SortOrder != "ASC" {
		t.Errorf(
			"expected sort order ASC, got %s",
			productRepository.receivedFilter.SortOrder,
		)
	}
}

func TestProductService_GetProducts_ReturnsErrorForInvalidStatus(t *testing.T) {
	productRepository := &fakeProductRepository{}
	categoryRepository := &fakeCategoryRepository{}

	service := NewProductService(productRepository, categoryRepository)

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

	if productRepository.getProductsCalls != 0 {
		t.Errorf(
			"expected repository not to be called, got %d calls",
			productRepository.getProductsCalls,
		)
	}
}

func TestProductService_GetProducts_CalculatesTotalPages(t *testing.T) {
	productRepository := &fakeProductRepository{totalItems: 21}
	categoryRepository := &fakeCategoryRepository{}

	service := NewProductService(productRepository, categoryRepository)

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

func TestProductService_GetProducts_CapsPageSizeAt100(t *testing.T) {
	productRepository := &fakeProductRepository{totalItems: 21}
	categoryRepository := &fakeCategoryRepository{}
	service := NewProductService(productRepository, categoryRepository)

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

	if productRepository.receivedFilter.Pagination.PageSize != 100 {
		t.Errorf(
			"expected repository to receive page size 100, got %d",
			productRepository.receivedFilter.Pagination.PageSize,
		)
	}

	if productList.Pagination.PageSize != 100 {
		t.Errorf(
			"expected response page size 100, got %d",
			productList.Pagination.PageSize,
		)
	}
}

func TestProductService_CreateProduct_DoesNotValidateCategoryWhenCategoryIsNil(t *testing.T) {
	productRepository := &fakeProductRepository{
		createdProduct: model.Product{
			ID:    1,
			Name:  "BANANA",
			Price: 30000,
		},
	}
	categoryRepository := &fakeCategoryRepository{}

	service := NewProductService(productRepository, categoryRepository)

	product := model.Product{
		Name:  "BANANA",
		Price: 30000,
	}

	createdProduct, err := service.CreateProduct(
		context.Background(),
		product,
	)

	if createdProduct.ID != 1 {
		t.Errorf("expected product ID 1, got %d", createdProduct.ID)
	}

	if categoryRepository.getCategoryByIDCalls != 0 {
		t.Errorf(
			"expected repository not to be called, got %d calls",
			categoryRepository.getCategoryByIDCalls,
		)
	}

	if productRepository.createProductCalls != 1 {
		t.Errorf(
			"expected repository to be called, got %d calls",
			productRepository.createProductCalls,
		)
	}

	if createdProduct.Category != nil {
		t.Errorf(
			"expected category to be nil",
		)
	}

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestProductService_CreateProduct_ReturnsErrorWhenCategoryDoesNotExist(t *testing.T) {
	productRepository := &fakeProductRepository{}
	categoryRepository := &fakeCategoryRepository{
		err: repository.ErrCategoryNotFound,
	}

	service := NewProductService(productRepository, categoryRepository)

	product := model.Product{
		Price: 30000,
		Name:  "BANANA",
		Category: &model.Category{
			ID: 2,
		},
	}

	product, err := service.CreateProduct(context.Background(), product)

	if categoryRepository.getCategoryByIDCalls != 1 {
		t.Errorf(
			"expected repository to be called, got %d calls",
			categoryRepository.getCategoryByIDCalls,
		)
	}

	if productRepository.createProductCalls != 0 {
		t.Errorf(
			"expected repository not to be called, got %d calls",
			productRepository.createProductCalls,
		)
	}

	if !errors.Is(err, repository.ErrCategoryNotFound) {
		t.Errorf(
			"expected err %s, got %s",
			repository.ErrCategoryNotFound, err.Error(),
		)
	}

}
