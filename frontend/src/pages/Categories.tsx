import { useCallback, useEffect, useState } from "react"
import { useAuth } from "../hooks/useAuth"
import { Alert, Box, Button, CircularProgress, Typography } from "@mui/material"
import { getCategories } from "../services/api"
import type { Category, CategoryPagination, CategoryQueryParams, CategorySort } from "../types/categories"
import CategoryTable from "../components/CategoryTable"
import CreateCategoryDialog from "../components/CreateCategoryDialog"
import UpdateCategoryDialog from "../components/UpdateCategoryDialog"
import DeleteCategoryDialog from "../components/DeleteCategoryDialog"
import { categorySortParams } from "../constants/categories"
import CustomPagination from "../components/CustomPagination"
import CategoryFilters from "../components/CategoryFilters"

export default function Categories() {
    const { token } = useAuth()

    const [error, setError] = useState<string | null>(null)

    const [categories, setCategories] = useState<Category[]>([])

    const [pagination, setPagination] = useState<CategoryPagination | null >(null)

    const [search, setSearch] = useState("")

    const [debouncedSearch, setDebouncedSearch] = useState("")

    const [page, setPage] = useState(1)

    const [categorySort, setCategorySort] = useState<CategorySort>("")

    const [isLoading, setIsLoading] = useState<boolean>(true)

    const [isFetching, setIsFetching] = useState<boolean>(false)

    const [openCreateDialog, setOpenCreateDialog] = useState(false)

    const [openUpdateDialog, setOpenUpdateDialog] = useState(false)
    
    const [openDeleteDialog, setOpenDeleteDialog] = useState(false)

    const [selectedCategory, setSelectedCategory] = useState<Category | null>(null)

    const pageSize = 10

    const handleOpenEditCategory = (category: Category) => {
        setSelectedCategory(category)
        setOpenUpdateDialog(true)
    }

    const handleOpenDeleteCategory = (category: Category) => {
        setSelectedCategory(category)
        setOpenDeleteDialog(true)
    }

    const handleCategoryDeleted = () => {
        if(pagination) {
            const newTotalItems = pagination?.totalItems - 1
            const newTotalPages = Math.max(1, Math.ceil(newTotalItems / pageSize))
            
            if (page > newTotalPages) {
                setPage(newTotalPages)
            } else {
                fetchCategories()
            }
        } else {
            setPage(1)
        }
    }
    
    const handleSearchChange = (searchStr: string) => {
        setSearch(searchStr)
        setPage(1)
    }

    const handleSortChange = (sort: CategorySort) => {
        setCategorySort(sort)
        setPage(1)
    }

    const handleResetFilters = () => {
        setSearch("")
        setDebouncedSearch("")
        setCategorySort("")
        setPage(1)
    }

    const fetchCategories = useCallback(async (signal?: AbortSignal) => {
        if (!token) {
            return
        }

        try {
            setIsFetching(true)
            setError(null)

            const {sortBy, sortOrder} = categorySortParams[categorySort]

            const params: CategoryQueryParams = {
                page,
                pageSize,
                search: debouncedSearch,
                sortBy,
                sortOrder
            }

            const response = await getCategories(token, params, signal)

            setCategories(response.data)
            setPagination(response.pagination)

        } catch (error) {
            if (error instanceof Error && error.name === "AbortError") {
                return
            }
            
            console.error(error)
            setError("Failed to load categories")
        } finally {
            if (!signal?.aborted) {
                setIsLoading(false)
                setIsFetching(false)
            }
        }
    }, [token, page, debouncedSearch, categorySort])


    useEffect(() => {
        const controller = new AbortController
        
        fetchCategories(controller.signal)

        return () => {
            controller.abort()
        };
    },[token, fetchCategories])

    useEffect(() => {
        const timer = setTimeout(()=>{
            setDebouncedSearch(search)
        }, 500)

        return () => {
            clearTimeout(timer)
        }
    }, [search])


    if (isLoading) {
        return (
            <Box
                sx={{
                    display: "flex",
                    justifyContent: "center",
                    alignItems: "center",
                    minHeight: "300px",
                }}
            >
                <CircularProgress />
            </Box>
        )
    }

    if (error) {
        return <Alert severity="error">{error}</Alert>
    }

    return(
        <Box>
            <Box
                sx={{
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "center",
                    mb: 2,
                }}
            >
                <Box
                    sx={{
                        display: "flex",
                        alignItems: "center",
                        gap: 1,
                    }}
                >
                    <Typography variant="h4">
                        Categories
                    </Typography>

                </Box>
                

                <Button
                    variant="contained"
                    onClick={() => {
                        setOpenCreateDialog(true)
                    }}
                >
                    Add Category
                </Button>
            </Box>

            <CategoryFilters
            search={search}
            categorySort={categorySort}
            onSortChange={handleSortChange}
            onSearchChange={handleSearchChange}
            onResetFilters={handleResetFilters}
            />

            <CategoryTable 
            categories={categories} 
            isFetching={isFetching} 
            onEdit={handleOpenEditCategory} 
            onDelete={handleOpenDeleteCategory}            
            />

            <CreateCategoryDialog
            open={openCreateDialog}
            onClose={() => setOpenCreateDialog(false)}
            onCreated={() => fetchCategories()}
            />

            {pagination &&
                <CustomPagination
                currentItemCount={categories.length}
                page={page}
                totalPages={pagination.totalPages}
                totalItems={pagination.totalItems}
                onPageChange={setPage}
                />
            }

            <UpdateCategoryDialog
            key={selectedCategory?.id ?? "update-category"}
            open={openUpdateDialog}
            category={selectedCategory}
            onClose={() => setOpenUpdateDialog(false)}
            onUpdated={() => fetchCategories()}
            />

            <DeleteCategoryDialog
            open={openDeleteDialog}
            category={selectedCategory}
            onClose={() => setOpenDeleteDialog(false)}
            onDeleted={handleCategoryDeleted}
            />
        </Box>
    )
}