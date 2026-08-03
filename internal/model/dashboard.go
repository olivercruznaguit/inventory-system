package model

type InventoryDashboard struct {
    TotalProducts       int
    ActiveProducts      int
    InactiveProducts    int
    TotalQuantityOnHand int
    LowStockProducts    int
    OutOfStockProducts  int
    TotalInventoryValue float64
}