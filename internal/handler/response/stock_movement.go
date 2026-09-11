package response

import (
	"time"

	"github.com/olivercruznaguit/inventory-system/internal/model"
)

type StockMovementResponse struct {
	ID                int       `json:"id"`
	Type              string    `json:"type"`
	Quantity          int       `json:"quantity"`
	RemainingQuantity int       `json:"remainingQuantity"`
	Reason            *string   `json:"reason"`
	CreatedAt         time.Time `json:"createdAt"`
}

func NewStockMovementResponse(stockMovement model.StockMovement) StockMovementResponse {
	return StockMovementResponse{
		ID:                int(stockMovement.ID),
		Type:              string(stockMovement.Type),
		Quantity:          stockMovement.Quantity,
		RemainingQuantity: stockMovement.RemainingQuantity,
		Reason:            stockMovement.Reason,
		CreatedAt:         stockMovement.CreatedAt,
	}
}

type StockMovementListResponse struct {
	Data []StockMovementResponse `json:"data"`
}

type RecentStockMovementResponse struct {
	ID                int       `json:"id"`
	Type              string    `json:"type"`
	Quantity          int       `json:"quantity"`
	ProductID         int       `json:"productId"`
	ProductName       string    `json:"productName"`
	RemainingQuantity int       `json:"remainingQuantity"`
	Reason            *string   `json:"reason"`
	CreatedAt         time.Time `json:"createdAt"`
}

func NewRecentStockMovementResponse(recentStockMovement model.RecentStockMovement) RecentStockMovementResponse {
	return RecentStockMovementResponse{
		ID:                int(recentStockMovement.ID),
		Type:              string(recentStockMovement.Type),
		ProductID:         int(recentStockMovement.ProductID),
		ProductName:       recentStockMovement.ProductName,
		Quantity:          recentStockMovement.Quantity,
		RemainingQuantity: recentStockMovement.RemainingQuantity,
		Reason:            recentStockMovement.Reason,
		CreatedAt:         recentStockMovement.CreatedAt,
	}
}

type RecentStockMovementListResponse struct {
	Data []RecentStockMovementResponse `json:"data"`
}
