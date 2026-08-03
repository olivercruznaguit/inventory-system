package response

import "github.com/olivercruznaguit/inventory-system/internal/model"

type InventoryDashboardResponse struct {
	TotalProducts       int     `json:"totalProducts"`
	ActiveProducts      int     `json:"activeProducts"`
	InactiveProducts    int     `json:"inactiveProducts"`
	TotalQuantityOnHand int     `json:"totalQuantityOnHand"`
	LowStockProducts    int     `json:"lowStockProducts"`
	OutOfStockProducts  int     `json:"outOfStockProducts"`
	TotalInventoryValue float64 `json:"totalInventoryValue"`
}

func NewInventoryDashboardResponse(dashboard model.InventoryDashboard) InventoryDashboardResponse {
	return InventoryDashboardResponse{
		TotalProducts:       dashboard.TotalProducts,
		ActiveProducts:      dashboard.ActiveProducts,
		InactiveProducts:    dashboard.InactiveProducts,
		TotalQuantityOnHand: dashboard.TotalQuantityOnHand,
		LowStockProducts:    dashboard.LowStockProducts,
		OutOfStockProducts:  dashboard.OutOfStockProducts,
		TotalInventoryValue: dashboard.TotalInventoryValue,
	}
}
