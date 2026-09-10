import { Button, Dialog, DialogActions, DialogContent, DialogTitle, TextField } from "@mui/material";
import { useState } from "react";
import { useAuth } from "../hooks/useAuth";
import { createCategory } from "../services/api";
import { useSnackbar } from "notistack";

type CreateCategoryDialogProps = {
    open: boolean
    onClose: () => void
    onCreated: () => void
}

export default function CreateCategoryDialog({ open, onClose, onCreated }: CreateCategoryDialogProps) {
    const { token } = useAuth()

    const { enqueueSnackbar } = useSnackbar()

    const [isCreating, setIsCreating] = useState(false)

    const [name, setName] = useState("")

    const [error, setError] = useState<string | null>(null)

    function handleClose() {
        if (isCreating) {
            return
        }

        setName("")
        setError(null)
        onClose()
    }

    async function handleSubmit() {
        if(!token) {
            return
        }

        if(!name.trim()) {
            setError("Name is required")
            return
        }

        try {
            setIsCreating(true)
            setError(null)

            await createCategory(token, name.trim())

            setName("")
            onClose()
            onCreated()

            enqueueSnackbar("Category added", {variant: "success"})
        }  catch (error) {
            console.error(error)

            if (error instanceof Error && error.cause === 409) {
                setError("A category with this name already exists")
                return
            }

            enqueueSnackbar("Failed to add category", { variant: "error" })
            setError("Failed to create category")
        } finally {
            setIsCreating(false)
        }
    }

    return (
        <Dialog 
        open={open}
        onClose={handleClose}
        maxWidth="sm"
        fullWidth>

            <DialogTitle>
                Create Category
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
                disabled={isCreating}>
                    {isCreating ? "Creating..." : "Create"}
                </Button>
            </DialogActions>
            
        </Dialog>
    )
}