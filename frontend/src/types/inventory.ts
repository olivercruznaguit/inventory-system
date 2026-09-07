export type InventoryRequest = {
    quantity: number
    reason: string
}

export type StockMovement = {
    id: number
    type: "IN" | "OUT"
    quantity: number
    remainingQuantity: number
    reason: string | null
    createdAt: string
}

export type StockMovementResponse = {
    data: StockMovement[]
}