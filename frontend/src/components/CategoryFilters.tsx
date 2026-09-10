import { Box, Button, FormControl, InputLabel, MenuItem, Select, TextField } from "@mui/material"
import type { CategorySort } from "../types/categories"
import { categorySortOptions } from "../constants/categories"

type CategoryFiltersProps = {
    search: string
    categorySort: CategorySort
    onSearchChange: (search: string) => void
    onSortChange: (sort: CategorySort) => void
    onResetFilters: () => void
}

export default function CategoryFilters({search, categorySort, onSearchChange, onSortChange, onResetFilters}: CategoryFiltersProps) {
    return (
        <Box
        sx={{
            display: "flex",
            gap: 1,
            mb: 2,
            flexWrap: "wrap",
        }}
        >
            <TextField 
            size="small"
            placeholder="Search"
            value={search} 
            sx={{ minWidth: 240 }}
            onChange={(event)=>{
                onSearchChange(event.target.value)
            }}/>

            <FormControl size="small" sx={{ minWidth: 180 }}>
                <InputLabel>Sort</InputLabel>

                <Select
                    value={categorySort}
                    label="Sort"
                    onChange={(event) => {
                        onSortChange(event.target.value as CategorySort)
                    }}
                >
                   {categorySortOptions.map((option)=>(
                        <MenuItem key={option.value} value={option.value}>
                            {option.label}
                        </MenuItem> 
                    ))}
                </Select>
            </FormControl>
            
            <Button
                size="small"
                variant="outlined"
                onClick={onResetFilters}
            >
                Reset Filters
            </Button>
        </Box>
    )
}