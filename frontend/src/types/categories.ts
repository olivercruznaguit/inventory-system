export type CategorySort =
    | ""
    | "name-asc"
    | "name-desc"

export type Category = {
    id: number
    name: string
    productCount: number
    createdAt: string
    updatedAt: string
}

export type CategoryPagination = {
    page: number
    pageSize: number
    totalItems: number
    totalPages: number
}

export type CategoryListResponse = {
    data: Category[]
    pagination: CategoryPagination
}

export type CategoryQueryParams = {
    page: number
    pageSize: number
    search: string
    sortBy: string
    sortOrder: string
}