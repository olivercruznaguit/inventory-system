import { Button, Dialog, DialogActions, DialogContent, DialogTitle, Typography } from "@mui/material"
import { useAuth } from "../hooks/useAuth"
import type { Category } from "../types/categories"
import { useState } from "react"
import { deleteCategory } from "../services/api"
import { useSnackbar } from "notistack"

type DeleteCategoryDialog = {
    open: boolean
    category: Category | null
    onClose: () => void
    onDeleted: () => void
}

export default function DeleteCategoryDialog({ open, category, onClose, onDeleted}: DeleteCategoryDialog) {
    const { token } = useAuth()

    const { enqueueSnackbar } = useSnackbar()

    const [isDeleting, setIsDeleting] = useState(false)

    async function handleSubmit() {
        if(!token || !category) {
            return
        }

        try {
            setIsDeleting(true)
            
            await deleteCategory(token, category.id)

            onClose()
            onDeleted()

            enqueueSnackbar("Category deleted", {variant: "success"})
        } catch (error) {
            console.error(error)

            enqueueSnackbar("Failed to delete category", { variant: "error" })
        } finally {
            setIsDeleting(false)
        }
    }


    function handleClose() {
        if (isDeleting) {
            return
        }

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
                <Typography variant="body1" sx={{ mb: 1 }}>
                    Are you sure you want to delete "{category?.name}"?
                </Typography>

                <Typography variant="body2" sx={{ color: "text.secondary" }}>
                    This action cannot be undone.
                </Typography>
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
                onClick={handleSubmit}
                disabled={isDeleting}
                >
                    { isDeleting ? "Deleting..." : "Delete"}
                </Button>
            </DialogActions>

        </Dialog>
    )
}