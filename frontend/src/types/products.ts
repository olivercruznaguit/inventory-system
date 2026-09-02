export type ProductStatus = "ACTIVE" | 'INACTIVE'

export type ProductSort =
    | ""
    | "name-asc"
    | "name-desc"
    | "price-asc"
    | "price-desc"

export type ProductCategory = {
    id: number
    name: string
}

export type Product = {
    id: number
    name: string
    price: number
    status: ProductStatus
    quantity: number
    minimumStock: number
    category?: ProductCategory
    createdAt: string
    updatedAt: string
}

export type ProductPagination = {
    page: number
    pageSize: number
    totalItems: number
    totalPages: number
}

export type ProductListResponse = {
    data: Product[]
    pagination: ProductPagination
}

export type CreateProductRequest = {
    name: string
    price: number
    categoryId?: number
    minimumStock: number
}

export type UpdateProductRequest = {
    name: string
    price: number
    status: ProductStatus
    categoryId?: number
    minimumStock: number
}

export type ProductQueryParams = {
    page: number
    pageSize: number
    status: ProductStatus
    search: string
    categoryId: string
    sortBy: string
    sortOrder: string
}