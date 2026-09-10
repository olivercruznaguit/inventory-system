import { Button, Dialog, DialogActions, DialogContent, DialogTitle, Typography } from "@mui/material"
import type { Product } from "../types/products"
import { useState } from "react"
import { useAuth } from "../hooks/useAuth"
import { deleteProduct } from "../services/api"
import { useSnackbar } from "notistack"

type DeleteProductDialogProps = {
    open: boolean
    product: Product | null
    onClose: () => void
    onDeleted: () => void
}

export default function DeleteProductDialog({ open, product, onDeleted, onClose }:DeleteProductDialogProps){
    const { token } = useAuth()

    const { enqueueSnackbar } = useSnackbar()
    
    const [isDeleting, setIsDeleting] = useState(false)

    async function handleSubmit() {
        if (!token || !product) {
            return
        }

        try{
            setIsDeleting(true)

            await deleteProduct(token, product.id)

            onClose()
            onDeleted()

            enqueueSnackbar("Product deleted", { variant: "success" })
        } catch (error) {
            console.error(error)

            enqueueSnackbar("Failed to delete product", { variant: "error" })
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

    return(
        <Dialog 
        open={open} 
        onClose={handleClose}>
            <DialogTitle>
                Delete Product
            </DialogTitle>

            <DialogContent>
                <Typography variant="body1" sx={{ mb: 1 }}>
                    Are you sure you want to delete "{product?.name}"?    
                </Typography>

                <Typography variant="body2" sx={{ color: "text.secondary" }}>
                    This action cannot be undone.   
                </Typography>
            </DialogContent>

            <DialogActions>
                <Button 
                onClick={handleClose}>
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