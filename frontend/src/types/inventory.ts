export type InventoryRequest = {
    quantity: number
    reason: string
}

export type StockMovementType = "IN" | "OUT"

export type StockMovement = {
    id: number
    type: StockMovementType
    quantity: number
    remainingQuantity: number
    reason: string | null
    createdAt: string
}

export type StockMovementResponse = {
    data: StockMovement[]
}