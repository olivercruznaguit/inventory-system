import { Box, Button, FormControl, InputLabel, MenuItem, Select, TextField } from "@mui/material"
import type { Category } from "../types/categories"
import type { ProductSort, ProductStatus } from "../types/products"
import { productSortOptions } from "../constants/products"

type ProductFiltersProps = {
    search: string
    status: ProductStatus
    categoryFilter: string
    productSort: ProductSort
    categories: Category[]
    onSearchChange: (search: string) => void
    onStatusChange: (status: ProductStatus) => void
    onCategoryChange: (categoryID: string) => void
    onSortChange: (sort: ProductSort) => void
    onResetFilters: () => void
}

export default function ProductFilters({
    search, 
    status, 
    categoryFilter, 
    productSort, 
    categories,
    onSearchChange,
    onStatusChange,
    onCategoryChange,
    onSortChange,
    onResetFilters
}:ProductFiltersProps) {
    return (
        <Box
        sx={{
            display: "flex",
            gap: 1,
            mb: 2
        }}
        >
            <TextField 
            size="small"
            placeholder="Search"
            value={search} 
            onChange={(event)=>{
                onSearchChange(event.target.value)
            }}/>

            <FormControl size="small" sx={{ minWidth: 180}}>
                <InputLabel>Status</InputLabel>

                <Select
                    value={status}
                    label="Status"
                    onChange={(event) => {
                        onStatusChange(event.target.value as ProductStatus)
                    }}
                >
                    <MenuItem value="ACTIVE">
                        Active
                    </MenuItem>

                    <MenuItem value="INACTIVE">
                        Inactive
                    </MenuItem>
                </Select>
            </FormControl>

            <FormControl size="small" sx={{ minWidth: 180}}>
                <InputLabel>Category</InputLabel>

                <Select
                    sx={{ maxHeight: "180px" }}
                    value={categoryFilter}
                    label="Category"
                    onChange={(event) => {
                        onCategoryChange(String(event.target.value))
                    }}
                >
                    <MenuItem value="">
                        All
                    </MenuItem>

                    {categories.map((category)=>(
                        <MenuItem key={category.id} value={category.id}>
                            {category.name}
                        </MenuItem>
                    ))}
                </Select>
            </FormControl>

            <FormControl size="small" sx={{ minWidth: 180 }}>
                <InputLabel>Sort</InputLabel>

                <Select
                    value={productSort}
                    label="Sort"
                    onChange={(event) => {
                        onSortChange(event.target.value as ProductSort)
                    }}
                >

                    {productSortOptions.map((option)=>(
                        <MenuItem key={option.value} value={option.value}>
                            {option.label}
                        </MenuItem>
                    ))}
                </Select>
            </FormControl>

            <Button
                size="small"
                variant="contained"
                onClick={onResetFilters}
                sx={{maxHeight: 56}}
            >
                Reset Filters
            </Button>
        </Box>
    )
}