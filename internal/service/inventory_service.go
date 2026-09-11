package service

import (
	"context"

	"github.com/olivercruznaguit/inventory-system/internal/database"
	"github.com/olivercruznaguit/inventory-system/internal/model"
	"github.com/olivercruznaguit/inventory-system/internal/repository"
)

type InventoryService struct {
	db *database.Database
}

func NewInventoryService(db *database.Database) *InventoryService {
	return &InventoryService{
		db: db,
	}
}

func (s *InventoryService) StockIn(ctx context.Context, request model.StockRequest) (model.Product, error) {
	if request.Quantity <= 0 {
		return model.Product{}, ErrInvalidProductQuantity
	}

	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return model.Product{}, err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	productRepo := repository.NewProductRepository(tx)
	movementRepo := repository.NewStockMovementRepository(tx)

	product, err := productRepo.GetProductByIDForUpdate(ctx, request.ProductID)
	if err != nil {
		return model.Product{}, err
	}

	newQuantity := product.Quantity + request.Quantity

	if err := productRepo.UpdateProductQuantity(ctx, request.ProductID, newQuantity); err != nil {
		return model.Product{}, err
	}

	movement := model.StockMovement{
		ProductID:         int64(request.ProductID),
		Type:              model.StockMovementTypeIn,
		Quantity:          request.Quantity,
		RemainingQuantity: newQuantity,
		Reason:            request.Reason,
	}

	if err := movementRepo.Create(ctx, &movement); err != nil {
		return model.Product{}, err
	}

	product.Quantity = newQuantity

	if err := tx.Commit(ctx); err != nil {
		return model.Product{}, err
	}

	return product, nil
}

func (s *InventoryService) StockOut(ctx context.Context, request model.StockRequest) (model.Product, error) {
	if request.Quantity <= 0 {
		return model.Product{}, ErrInvalidProductQuantity
	}

	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return model.Product{}, err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	productRepo := repository.NewProductRepository(tx)
	movementRepo := repository.NewStockMovementRepository(tx)

	product, err := productRepo.GetProductByIDForUpdate(ctx, request.ProductID)
	if err != nil {
		return model.Product{}, err
	}

	if product.Quantity < request.Quantity {
		return model.Product{}, ErrInsufficientStock
	}

	newQuantity := product.Quantity - request.Quantity

	if err := productRepo.UpdateProductQuantity(ctx, request.ProductID, newQuantity); err != nil {
		return model.Product{}, err
	}

	movement := model.StockMovement{
		ProductID:         int64(request.ProductID),
		Type:              model.StockMovementTypeOut,
		Quantity:          request.Quantity,
		RemainingQuantity: newQuantity,
		Reason:            request.Reason,
	}

	if err := movementRepo.Create(ctx, &movement); err != nil {
		return model.Product{}, err
	}

	product.Quantity = newQuantity

	if err := tx.Commit(ctx); err != nil {
		return model.Product{}, err
	}

	return product, nil
}

func (s *InventoryService) GetStockMovements(ctx context.Context, productID int) ([]model.StockMovement, error) {
	productRepo := repository.NewProductRepository(s.db.DB())
	movementRepo := repository.NewStockMovementRepository(s.db.DB())

	_, err := productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return []model.StockMovement{}, err
	}

	return movementRepo.GetByProductID(ctx, productID)
}

func (s *InventoryService) GetInventoryDashboard(ctx context.Context) (model.InventoryDashboard, error) {
	productRepo := repository.NewProductRepository(s.db.DB())

	return productRepo.GetInventoryDashboard(ctx)
}

func (sr *InventoryService) GetRecentStockMovements(ctx context.Context, limit int) ([]model.RecentStockMovement, error) {
	if limit <= 0 {
		limit = 10
	}

	inventoryRepo := repository.NewStockMovementRepository(sr.db.DB())

	recentStockMovements, err := inventoryRepo.GetRecent(ctx, limit)

	if err != nil {
		return []model.RecentStockMovement{}, err
	}

	return recentStockMovements, nil
}
