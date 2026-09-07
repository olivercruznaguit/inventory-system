import { Box, Button, Dialog, DialogActions, DialogContent, DialogTitle, TextField, Typography } from "@mui/material"
import type { Product } from "../types/products"
import { useState } from "react"
import type { InventoryRequest } from "../types/inventory"
import { useAuth } from "../hooks/useAuth"
import { stockIn } from "../services/api"
import { useSnackbar } from "notistack"

type ProductStockInDialogProp = {
    product: Product | null
    open: boolean
    onClose: () => void
    onSubmit: () => void
}

export default function ProductStockInDialog({product, open, onClose, onSubmit}: ProductStockInDialogProp) {
    const { enqueueSnackbar } = useSnackbar();
    
    const { token } = useAuth()

    const [error, setError] = useState<string | null>(null)
    
    const [isSubmitting, setIsSubmitting] = useState<boolean>(false)

    const [form, setForm] = useState<InventoryRequest>({
        quantity: 1,
        reason: ""
    })

    async function handleSubmit() {
        if (!token || !product) {
            return
        }

        if (form.quantity <= 0) {
            setError("Invalid quantity")
            return
        }

        try{
            setError(null)
            setIsSubmitting(true)

            await stockIn(token, product.id, form)

            onClose()
            onSubmit()
            enqueueSnackbar(
                `Successfully added ${form.quantity} units`,
                { variant: "success" }
            )
        } catch (error) {
            console.error(error)
            setError("Failed to stock in product")
        } finally {
            setIsSubmitting(false)
        }
    }

    return (
        <Dialog
        open={open}
        onClose={onClose}
        fullWidth
        maxWidth="sm"
        >
            <DialogTitle>
                Stock In
            </DialogTitle>

            <DialogContent>
                <Box 
                sx={{ 
                    display: "flex",
                    flexDirection: "column",
                    gap: 1
                }}
                >
                    <Box>
                        <Typography variant="subtitle1">
                            Product
                        </Typography>

                        <Typography variant="subtitle1" sx={{fontWeight: "bold"}}>
                            {product?.name}
                        </Typography>
                    </Box>

                    <Box>
                        <Typography variant="subtitle1">
                            Current Stock
                        </Typography>

                        <Typography variant="subtitle1" sx={{fontWeight: "bold"}}>
                            {product?.quantity}
                        </Typography>
                    </Box>
                </Box>

                <TextField
                value={form.quantity}
                label="Quantity"
                margin="normal"
                fullWidth
                type="number"
                onChange={(event) =>
                        setForm({
                            ...form,
                            quantity: Number(event.target.value),
                        })
                    }
                />

                 <TextField
                label="Reason (optional)"
                margin="normal"
                fullWidth
                onChange={(event) =>
                    setForm({
                        ...form,
                        reason: event.target.value,
                    })
                }
                />

                { error && (
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
                onClick={onClose}
                >
                    Cancel
                </Button>

                <Button
                variant="contained"
                onClick={handleSubmit}
                disabled={isSubmitting}
                >
                    { isSubmitting ? "Stocking in..." : "Stock In" }
                </Button>
            </DialogActions>

        </Dialog>
    )
}