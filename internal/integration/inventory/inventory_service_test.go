package inventory_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/olivercruznaguit/inventory-system/internal/integration/testutil"
	"github.com/olivercruznaguit/inventory-system/internal/model"
	"github.com/olivercruznaguit/inventory-system/internal/repository"
	"github.com/olivercruznaguit/inventory-system/internal/service"
)

func TestInventoryService_StockIn_Success(t *testing.T) {
	// ARRANGE
	db := testutil.SetupTestDatabase(t)

	ctx := context.Background()

	testutil.CleanupDatabase(t, ctx, db)

	category := testutil.SeedCategory(t, ctx, db,
		model.Category{
			Name: "Electronics",
		},
	)

	product := testutil.SeedProduct(t, ctx, db, model.Product{
		Name:     "Laptop",
		Price:    1000.00,
		Status:   model.ProductStatusActive,
		Quantity: 100,
		Category: &category,
	})

	inventoryService := service.NewInventoryService(db)

	request := model.StockRequest{
		ProductID: int(product.ID),
		Quantity:  10,
	}

	// ACT
	updatedProduct, err := inventoryService.StockIn(ctx, request)

	// ASSERT
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updatedProduct.Quantity != 110 {
		t.Errorf("expected quantity 110, got %d", updatedProduct.Quantity)
	}

	testutil.AssertProductQuantity(t, ctx, db, product.ID, 110)

	testutil.AssertStockMovementCount(t, ctx, db, product.ID, 1)

	testutil.AssertStockMovement(t, ctx, db, product.ID, model.StockMovement{
		Type:              model.StockMovementTypeIn,
		Quantity:          10,
		RemainingQuantity: 110,
	})
}

func TestInventoryService_StockIn_InvalidQuantity(t *testing.T) {
	// ARRANGE
	db := testutil.SetupTestDatabase(t)

	ctx := context.Background()

	testutil.CleanupDatabase(t, ctx, db)

	category := testutil.SeedCategory(t, ctx, db,
		model.Category{
			Name: "Electronics",
		},
	)

	product := testutil.SeedProduct(t, ctx, db, model.Product{
		Name:     "Laptop",
		Price:    1000.00,
		Status:   model.ProductStatusActive,
		Quantity: 100,
		Category: &category,
	})

	inventoryService := service.NewInventoryService(db)

	request := model.StockRequest{
		ProductID: int(product.ID),
		Quantity:  0,
	}

	// ACT
	_, err := inventoryService.StockIn(ctx, request)

	// ASSERT
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	testutil.AssertProductQuantity(t, ctx, db, product.ID, 100)

	if !errors.Is(err, service.ErrInvalidProductQuantity) {
		t.Fatalf("expected error %v, got %v", service.ErrInvalidProductQuantity, err)
	}

	testutil.AssertStockMovementCount(t, ctx, db, product.ID, 0)
}

func TestInventoryService_StockIn_ProductNotFound(t *testing.T) {
	// ARRANGE
	db := testutil.SetupTestDatabase(t)

	ctx := context.Background()

	testutil.CleanupDatabase(t, ctx, db)

	inventoryService := service.NewInventoryService(db)

	request := model.StockRequest{
		ProductID: 999, // Non-existent product ID
		Quantity:  10,
	}

	// ACT
	_, err := inventoryService.StockIn(ctx, request)

	// ASSERT
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, repository.ErrProductNotFound) {
		t.Fatalf("expected error %v, got %v", repository.ErrProductNotFound, err)
	}
}

func TestInventoryService_StockOut_Success(t *testing.T) {
	// ARRANGE
	db := testutil.SetupTestDatabase(t)

	ctx := context.Background()

	testutil.CleanupDatabase(t, ctx, db)

	category := testutil.SeedCategory(t, ctx, db,
		model.Category{
			Name: "Electronics",
		},
	)

	product := testutil.SeedProduct(t, ctx, db, model.Product{
		Name:     "Laptop",
		Price:    1000.00,
		Status:   model.ProductStatusActive,
		Quantity: 110,
		Category: &category,
	})

	inventoryService := service.NewInventoryService(db)

	request := model.StockRequest{
		ProductID: int(product.ID),
		Quantity:  10,
	}

	// ACT
	updatedProduct, err := inventoryService.StockOut(ctx, request)

	// ASSERT
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updatedProduct.Quantity != 100 {
		t.Errorf("expected quantity 100, got %d", updatedProduct.Quantity)
	}

	testutil.AssertProductQuantity(t, ctx, db, product.ID, 100)

	testutil.AssertStockMovementCount(t, ctx, db, product.ID, 1)

	testutil.AssertStockMovement(t, ctx, db, product.ID, model.StockMovement{
		Type:              model.StockMovementTypeOut,
		Quantity:          10,
		RemainingQuantity: 100,
	})

}

func TestInventoryService_StockOut_InvalidQuantity(t *testing.T) {
	// ARRANGE
	db := testutil.SetupTestDatabase(t)

	ctx := context.Background()

	testutil.CleanupDatabase(t, ctx, db)

	category := testutil.SeedCategory(t, ctx, db,
		model.Category{
			Name: "Electronics",
		},
	)

	product := testutil.SeedProduct(t, ctx, db, model.Product{
		Name:     "Laptop",
		Price:    1000.00,
		Status:   model.ProductStatusActive,
		Quantity: 100,
		Category: &category,
	})

	inventoryService := service.NewInventoryService(db)

	request := model.StockRequest{
		ProductID: int(product.ID),
		Quantity:  0,
	}

	// ACT
	_, err := inventoryService.StockOut(ctx, request)

	// ASSERT
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, service.ErrInvalidProductQuantity) {
		t.Fatalf("expected error %v, got %v", service.ErrInvalidProductQuantity, err)
	}

	testutil.AssertProductQuantity(t, ctx, db, product.ID, 100)

	testutil.AssertStockMovementCount(t, ctx, db, product.ID, 0)
}

func TestInventoryService_StockOut_InsufficientStock(t *testing.T) {
	// ARRANGE
	db := testutil.SetupTestDatabase(t)

	ctx := context.Background()

	testutil.CleanupDatabase(t, ctx, db)

	category := testutil.SeedCategory(t, ctx, db,
		model.Category{
			Name: "Electronics",
		},
	)

	product := testutil.SeedProduct(t, ctx, db, model.Product{
		Name:     "Laptop",
		Price:    1000.00,
		Status:   model.ProductStatusActive,
		Quantity: 100,
		Category: &category,
	})

	inventoryService := service.NewInventoryService(db)

	request := model.StockRequest{
		ProductID: int(product.ID),
		Quantity:  110,
	}

	// ACT
	_, err := inventoryService.StockOut(ctx, request)

	// ASSERT
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, service.ErrInsufficientStock) {
		t.Fatalf("expected error %v, got %v", service.ErrInsufficientStock, err)
	}

	testutil.AssertProductQuantity(t, ctx, db, product.ID, 100)

	testutil.AssertStockMovementCount(t, ctx, db, product.ID, 0)
}

// Verifies that concurrent StockOut operations are serialized
// using PostgreSQL row-level locking (SELECT ... FOR UPDATE).
// The final quantity should be 0 with exactly 10 stock movements.
func TestInventoryService_StockOut_Concurrent(t *testing.T) {
	// ARRANGE
	db := testutil.SetupTestDatabase(t)

	ctx := context.Background()

	testutil.CleanupDatabase(t, ctx, db)

	category := testutil.SeedCategory(t, ctx, db,
		model.Category{
			Name: "Electronics",
		},
	)

	product := testutil.SeedProduct(t, ctx, db, model.Product{
		Name:     "Laptop",
		Price:    1000.00,
		Status:   model.ProductStatusActive,
		Quantity: 100,
		Category: &category,
	})

	inventoryService := service.NewInventoryService(db)

	request := model.StockRequest{
		ProductID: int(product.ID),
		Quantity:  10,
	}

	// ACT
	requests := 10
	ch := make(chan error, requests)
	var wg sync.WaitGroup

	for range requests {
		wg.Go(func() {
			_, err := inventoryService.StockOut(ctx, request)
			ch <- err
		})
	}

	wg.Wait()
	close(ch)

	// ASSERT
	for err := range ch {
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	}

	testutil.AssertProductQuantity(t, ctx, db, product.ID, 0)

	testutil.AssertStockMovementCount(t, ctx, db, product.ID, requests)
}
