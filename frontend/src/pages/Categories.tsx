import { useCallback, useEffect, useState } from "react"
import { useAuth } from "../hooks/useAuth"
import { Alert, Box, Button, CircularProgress, Typography } from "@mui/material"
import { getCategories } from "../services/api"
import type { Category } from "../types/categories"
import CategoryTable from "../components/CategoryTable"
import CreateCategoryDialog from "../components/CreateCategoryDialog"
import UpdateCategoryDialog from "../components/UpdateCategoryDialog"
import DeleteCategoryDialog from "../components/DeleteCategoryDialog"

export default function Categories() {
    const { token } = useAuth()

    const [error, setError] = useState<string | null>(null)

    const [categories, setCategories] = useState<Category[]>([])

    const [isLoading, setIsLoading] = useState<boolean>(true)

    const [isFetching, setIsFetching] = useState<boolean>(false)

    const [openCreateDialog, setOpenCreateDialog] = useState(false)

    const [openUpdateDialog, setOpenUpdateDialog] = useState(false)
    
    const [openDeleteDialog, setOpenDeleteDialog] = useState(false)

    const [selectedCategory, setSelectedCategory] = useState<Category | null>(null)

    const handleOpenEditCategory = (category: Category) => {
        setSelectedCategory(category)
        setOpenUpdateDialog(true)
    }

    const handleOpenDeleteCategory = (category: Category) => {
        setSelectedCategory(category)
        setOpenDeleteDialog(true)
    }

    const fetchCategories = useCallback(async (signal?: AbortSignal) => {
        if (!token) {
            return
        }

        try {
            setIsFetching(true)
            setError(null)

            const response = await getCategories(token, signal)

            setCategories(response)

        } catch (error) {
            console.error(error)
            setError("Failed to load categories")
        } finally {
            if (!signal?.aborted) {
                setIsLoading(false)
                setIsFetching(false)
            }
        }
    }, [token])


    useEffect(() => {
        const controller = new AbortController
        
        fetchCategories()

        return () => {
            controller.abort()
        };
    },[token, fetchCategories])

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
            onDeleted={() => fetchCategories()}
            />
        </Box>
    )
}