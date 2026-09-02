import type { ProductSort } from "../types/products"

export const productSortOptions = [
    { value: "", label: "Default" },
    { value: "name-asc", label: "Name: A → Z" },
    { value: "name-desc", label: "Name: Z → A" },
    { value: "price-asc", label: "Price: Low → High" },
    { value: "price-desc", label: "Price: High → Low" },
]

export const productSortParams: Record<ProductSort, { sortBy: string; sortOrder: string }> = {
    "": {
        sortBy: "",
        sortOrder: "",
    },
    "name-asc": {
        sortBy: "name",
        sortOrder: "ASC",
    },
    "name-desc": {
        sortBy: "name",
        sortOrder: "DESC",
    },
    "price-asc": {
        sortBy: "price",
        sortOrder: "ASC",
    },
    "price-desc": {
        sortBy: "price",
        sortOrder: "DESC",
    },
}