import type { AuthResponse } from "../types/auth"
import type { CategoryListResponse, CategoryQueryParams } from "../types/categories"
import type { InventoryDashboard } from "../types/dashboard"
import type { InventoryRequest, StockMovementResponse } from "../types/inventory"
import type { CreateProductRequest, ProductListResponse, ProductQueryParams, UpdateProductRequest } from "../types/products"

const API_URL = "http://localhost:8080"

export async function login(email: string, password: string):Promise<AuthResponse> {
    const response = await fetch(`${API_URL}/auth/login`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            email,
            password,
        })
    })

    if (!response.ok){
        throw new Error("Login failed")
    }

    return await response.json() 
}

// INVENTORY
export async function getInventoryDashboard(token: string):Promise<InventoryDashboard> {
    const response = await fetch(`${API_URL}/inventory/dashboard`, {
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
    })

    if (!response.ok){
        throw new Error("Failed to fetch inventory dashboard")
    }
    return response.json()
}

export async function stockIn(token: string, productID: string | number, request: InventoryRequest):Promise<void>{
    const response = await fetch(`${API_URL}/products/${productID}/stock-in`, 
        {
            method: "POST",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json",
            },
            body: JSON.stringify(request),
        }
    )

    if(!response.ok) {
        throw new Error("failed to stock in product")
    }
}

export async function stockOut(token: string, productID: string | number, request: InventoryRequest):Promise<void>{
    const response = await fetch(`${API_URL}/products/${productID}/stock-out`, 
        {
            method: "POST",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json",
            },
            body: JSON.stringify(request),
        }
    )

    if(!response.ok) {
        throw new Error("failed to stock out product")
    }
}

export async function stockMovement(token: string, productID: string | number): Promise<StockMovementResponse> {
    const response = await fetch(`${API_URL}/products/${productID}/stock-movements`, {
        method: "GET",
        headers: {
            "Authorization": `Bearer ${token}`,
            "Content-Type": "application/json"
        }
    })

    if (!response.ok) {
        throw new Error("Failed to get product stock movement")
    }

    return response.json()
}

//  PRODUCTS
export async function getProducts(token: string, params: ProductQueryParams, signal?: AbortSignal): Promise<ProductListResponse> {
    const urlParams = new URLSearchParams({
        page: params.page.toString(),
        pageSize: params.pageSize.toString(),
    })

    if (params.status) {
        urlParams.set("status", params.status)
    }

    if (params.search) {
        urlParams.set("search", params.search)
    }

    if(params.categoryId) {
        urlParams.set("categoryId", params.categoryId)
    }

    if(params.sortBy) {
        urlParams.set("sortBy", params.sortBy)
    }

    if(params.sortOrder) {
        urlParams.set("sortOrder", params.sortOrder)
    }
    
    const response = await fetch(
        `${API_URL}/products?${urlParams.toString()}`,
        {
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json",
            },
            signal
        }
    )

    if (!response.ok) {
        throw new Error("Failed to fetch products")
    }

    return response.json()
}

export async function createProduct(token: string, request: CreateProductRequest):Promise<void> {
    const response = await fetch(`${API_URL}/products`,{
        method: "POST",
        headers: {
            "Authorization": `Bearer ${token}`,
            "Content-Type": "application/json",
        },
        body: JSON.stringify(request),
    })

    if(!response.ok){
        throw new Error("Failed to create product")
    }
}

export async function updateProduct(token: string, productID: string | number, request: UpdateProductRequest):Promise<void> {
    const response = await fetch(`${API_URL}/products/${productID}`,{
        method: "PUT",
        headers: {
            "Authorization": `Bearer ${token}`,
            "Content-Type": "application/json",
        },
        body: JSON.stringify(request),
    })

    if(!response.ok){
        throw new Error("Failed to update product")
    }
}

export async function deleteProduct(token: string, productID: string | number):Promise<void> {
    const response = await fetch(`${API_URL}/products/${productID}`,{
        method: "DELETE",
        headers: {
            "Authorization": `Bearer ${token}`,
            "Content-Type": "application/json",
        },
    })

    if(!response.ok){
        throw new Error("Failed to delete product")
    }
}

// CATEGORIES

export async function getCategories(token: string, params: CategoryQueryParams, signal?: AbortSignal): Promise<CategoryListResponse> {
    const urlParams = new URLSearchParams({
        page: params.page.toString(),
        pageSize: params.pageSize.toString(),
    })  

    if (params.search) {
        urlParams.set("search", params.search)
    }

    if(params.sortBy) {
        urlParams.set("sortBy", params.sortBy)
    }

    if(params.sortOrder) {
        urlParams.set("sortOrder", params.sortOrder)
    }
    
    const response = await fetch(`${API_URL}/categories?${urlParams.toString()}`,
        {
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json",
            },
            signal
        }
    )

    if (!response.ok) {
        throw new Error("Failed to fetch categories")
    }

    return response.json()
}

export async function createCategory(token: string, name: string): Promise<void> {
    const response = await fetch(`${API_URL}/categories`,
        {
            method: "POST",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ name })
        }
    )

    if (!response.ok) {
        throw new Error("Failed to create category")
    }
}

export async function deleteCategory(token: string, categoryID: string | number): Promise<void> {
    const response = await fetch(`${API_URL}/categories/${categoryID}`,
        {
            method: "DELETE",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json",
            }
        }
    )

    if (!response.ok) {
        throw new Error("Failed to delete category")
    }
}

export async function updateCategory(token: string, categoryID: string | number, name: string): Promise<void> {
    const response = await fetch(`${API_URL}/categories/${categoryID}`,
        {
            method: "PUT",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ name })
        }
    )

    if (!response.ok) {
        throw new Error("Failed to update category")
    }
}