import type { StockMovementType } from "./inventory"

export type AlertSeverity = "success" | "warning" | "error"

export type InventoryDashboard = {
    totalProducts: number
    activeProducts: number
    inactiveProducts: number
    totalQuantityOnHand: number
    lowStockProducts: number
    outOfStockProducts: number
    totalInventoryValue: number
}

export type RecentStockMovement = {
    id: number
    type: StockMovementType
    productName: string
    productId: number
    quantity: number
    remainingQuantity: number
    reason: string | null
    createdAt: string
}

export type RecentStockMovementResponse = {
    data: RecentStockMovement[]
}
