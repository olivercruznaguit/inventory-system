import type { CategorySort } from "../types/categories";

export const categorySortOptions = [
    { value: "", label: "Default" },
    { value: "name-asc", label: "Name: A → Z" },
    { value: "name-desc", label: "Name: Z → A" },
]

export const categorySortParams: Record<CategorySort, { sortBy: string; sortOrder: string }> = {
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
    }
}