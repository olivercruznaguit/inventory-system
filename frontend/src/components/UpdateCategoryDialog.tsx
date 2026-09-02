import { Button, Dialog, DialogActions, DialogContent, DialogTitle, TextField, Typography } from "@mui/material"
import { useAuth } from "../hooks/useAuth"
import type { Category } from "../types/categories"
import { useState } from "react"
import { updateCategory } from "../services/api"

type UpdateCategoryDialogProps = {
    open: boolean
    category: Category | null
    onClose: () => void
    onUpdated: () => void
}


export default function UpdateCategoryDialog({ open, category, onClose, onUpdated }: UpdateCategoryDialogProps) {
    const { token } = useAuth() 

    const [isUpdating, setIsUpdating] = useState(false)

    const [error, setError] = useState<string | null>(null)

    const [name, setName] = useState(category?.name ?? "")

    async function handleSubmit() { 
        if(!token || !category) {
            setError("Invalid category or token")
            return
        }

        if(!name.trim()) {
            setError("Name is required")
            return
        }

        try {
            setIsUpdating(true)
            setError(null)


            await updateCategory(token, category.id, name)

            onClose()
            onUpdated()
        } catch (error) {
            console.error(error)
            setError("Failed to update category")
        } finally {
            setIsUpdating(false)
        }
    }

    // useEffect(()=>{
    //     if(category) {
    //         setName(category.name)
    //     } else {
    //         setName(null)
    //     }
    // },[category])


    // if(!name) {
    //     return null
    // }

    return (
        <Dialog 
        open={open}
        onClose={onClose}
        maxWidth="sm"
        fullWidth>

            <DialogTitle>
                Update Category
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
                disabled={isUpdating}>
                    {isUpdating ? "Updating..." : "Update"}
                </Button>
            </DialogActions>
            
        </Dialog>
    )

}