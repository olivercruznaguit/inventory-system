import { Button, Dialog, DialogActions, DialogContent, DialogTitle, Typography } from "@mui/material"
import { useAuth } from "../hooks/useAuth"
import type { Category } from "../types/categories"
import { useState } from "react"
import { deleteCategory } from "../services/api"

type DeleteCategoryDialog = {
    open: boolean
    category: Category | null
    onClose: () => void
    onDeleted: () => void
}

export default function DeleteCategoryDialog({ open, category, onClose, onDeleted}: DeleteCategoryDialog) {
    const { token } = useAuth()

    const [isDeleting, setIsDeleting] = useState(false)

    const [error, setError] = useState<string | null>(null)

    async function handleSubmit() {
        if(!token || !category) {
            return
        }

        try {
            setIsDeleting(true)
            
            await deleteCategory(token, category.id)

            handleClose()
            onDeleted()
        } catch (error) {
            console.error(error)
            setError("Failed to delete category")
        } finally {
            setIsDeleting(false)
        }
    }


    const handleClose = () => {
        setError(null)
        onClose()
    }

    return (
        <Dialog 
        open={open}
        onClose={handleClose}
        >
            <DialogTitle>
                Delete Category
            </DialogTitle>

            <DialogContent>
                <Typography variant="body1">
                    Are you sure you want to delete "{category?.name}"?    
                </Typography>

                <Typography variant="body1">
                    This action cannot be undone.   
                </Typography>

                {error && (
                    <Typography
                        color="error"
                        sx={{ mt: 2 }}
                    >
                        {error}
                    </Typography>
                )}
            </DialogContent>

            <DialogActions>
                <Button
                onClick={handleClose}
                >
                    Cancel
                </Button>

                <Button 
                variant="contained"
                color="error"
                sx={{ ml: 1 }}
                onClick={handleSubmit}
                >
                    { isDeleting ? "Deleting..." : "Delete"}
                </Button>
            </DialogActions>

        </Dialog>
    )
}