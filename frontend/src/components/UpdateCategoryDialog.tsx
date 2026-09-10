import { Button, Dialog, DialogActions, DialogContent, DialogTitle, TextField } from "@mui/material"
import { useAuth } from "../hooks/useAuth"
import type { Category } from "../types/categories"
import { useState } from "react"
import { updateCategory } from "../services/api"
import { useSnackbar } from "notistack"

type UpdateCategoryDialogProps = {
    open: boolean
    category: Category | null
    onClose: () => void
    onUpdated: () => void
}

export default function UpdateCategoryDialog({ open, category, onClose, onUpdated }: UpdateCategoryDialogProps) {
    const { token } = useAuth() 

    const { enqueueSnackbar } = useSnackbar()

    const [isUpdating, setIsUpdating] = useState(false)

    const [error, setError] = useState<string | null>(null)

    const [name, setName] = useState(category?.name ?? "")

    async function handleSubmit() { 
        if(!token || !category) {
            return
        }

        if(!name.trim()) {
            setError("Name is required")
            return
        }

        try {
            setIsUpdating(true)
            setError(null)

            await updateCategory(token, category.id, name.trim())
            
            onClose()
            onUpdated()

            enqueueSnackbar("Category updated", {variant: "success"})
        } catch (error) {
            console.error(error)

            if (error instanceof Error && error.cause === 409) {
                setError("A category with this name already exists")
                return
            }

            enqueueSnackbar("Failed to update category", { variant: "error" })
            setError("Failed to update category")
        } finally {
            setIsUpdating(false)
        }
    }

    function handleClose() {
        if (isUpdating) {
            return
        }

        setName("")
        setError(null)
        onClose()
    }
    
    return (
        <Dialog 
        open={open}
        onClose={handleClose}
        maxWidth="sm"
        fullWidth>

            <DialogTitle>
                Update Category
            </DialogTitle>

            <DialogContent>
                <TextField
                    fullWidth
                    required
                    label="Name"
                    margin="normal"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    error={Boolean(error)}
                    size="small"
                    helperText={error ?? "Enter a category name"}
                />
            </DialogContent>

            <DialogActions>
                <Button onClick={handleClose}>
                    Cancel
                </Button>

                <Button 
                variant="contained" 
                onClick={handleSubmit} 
                disabled={isUpdating}>
                    {isUpdating ? "Updating..." : "Update"}
                </Button>
            </DialogActions>
            
        </Dialog>
    )

}