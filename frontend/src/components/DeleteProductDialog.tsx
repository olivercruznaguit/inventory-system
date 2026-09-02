import { Button, Dialog, DialogActions, DialogContent, DialogTitle, Typography } from "@mui/material"
import type { Product } from "../types/products"
import { useState } from "react"
import { useAuth } from "../hooks/useAuth"
import { deleteProduct } from "../services/api"

type DeleteProductDialogProps = {
    open: boolean
    product: Product | null
    onClose: () => void
    onDeleted: () => void
}

export default function DeleteProductDialog({ open, product, onDeleted, onClose }:DeleteProductDialogProps){
    const { token } = useAuth()
    
    const [isDeleting, setIsDeleting] = useState<boolean>(false)

    const [error, setError] = useState<string | null>(null)

    async function handleSubmit() {
        if (!token || !product) {
            return
        }

        try{
            setIsDeleting(true)

            await deleteProduct(token, product.id)

            onClose()
            onDeleted()
        } catch (error) {
            console.error(error)
            setError("Failed to delete product")
        } finally {
            setIsDeleting(false)
        }
    }

    return(
        <Dialog 
        open={open} 
        onClose={onClose}>
            <DialogTitle>
                Delete Product
            </DialogTitle>

            <DialogContent>
                <Typography variant="body1">
                    Are you sure you want to delete "{product?.name}"?    
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
                onClick={onClose}>
                    Cancel
                </Button>

                <Button
                variant="contained"
                disabled={isDeleting}
                color="error"
                onClick={handleSubmit}
                >
                    {isDeleting ? "Deleting..." : "Delete"}    
                </Button>
            </DialogActions>
        </Dialog>
    )
}