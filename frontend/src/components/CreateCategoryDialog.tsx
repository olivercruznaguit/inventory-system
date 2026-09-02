import { Button, Dialog, DialogActions, DialogContent, DialogTitle, TextField, Typography } from "@mui/material";
import { useState } from "react";
import { useAuth } from "../hooks/useAuth";
import { createCategory } from "../services/api";

type CreateCategoryDialogProps = {
    open: boolean
    onClose: () => void
    onCreated: () => void
}

export default function CreateCategoryDialog({ open, onClose, onCreated }: CreateCategoryDialogProps) {
    const { token } = useAuth()

    const [isCreating, setIsCreating] = useState(false)

    const [name, setName] = useState("")

    const [error, setError] = useState<string | null>(null)

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

            await createCategory(token, name)

            onClose()
            onCreated()
        } catch (error) {
            console.error(error)
            setError("Failed to create category")
        } finally {
            setIsCreating(false)
        }
    }

    return (
        <Dialog 
        open={open}
        onClose={onClose}
        maxWidth="sm"
        fullWidth>

            <DialogTitle>
                Create Product
            </DialogTitle>

            <DialogContent>
                <TextField
                    fullWidth
                    label="Name"
                    margin="normal"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                />


                { error && (
                    <Typography 
                    color="error" 
                    sx={{ mt: 2 }}>
                        {error}
                    </Typography>
                )}

            </DialogContent>

            <DialogActions>
                <Button onClick={onClose}>
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