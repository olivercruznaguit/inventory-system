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

type FormErrors = {
    quantity?: string
}

export default function ProductStockInDialog({product, open, onClose, onSubmit}: ProductStockInDialogProp) {
    const { token } = useAuth()
    
    const { enqueueSnackbar } = useSnackbar();

    const [errors, setErrors] = useState<FormErrors>({})
    
    const [isSubmitting, setIsSubmitting] = useState(false)

    const [form, setForm] = useState<InventoryRequest>({
        quantity: 1,
        reason: ""
    })

    async function handleSubmit() {
        if (!token || !product) {
            return
        }

        const validationErrors = validateForm()

        if (Object.keys(validationErrors).length > 0) {
            setErrors(validationErrors)
            return
        }

        const quantity = form.quantity

        try{
            setErrors({})
            setIsSubmitting(true)

            await stockIn(token, product.id, form)

            setForm({
                quantity: 1,
                reason: ""
            })

            onClose()
            onSubmit()

            enqueueSnackbar(
                `Successfully added ${quantity} ${quantity > 1 ? "units" : "unit" }`,
                { variant: "success" }
            )
        } catch (error) {
            console.error(error)
            
            if (error instanceof Error) {
                enqueueSnackbar(error.message, { variant: "error" })
            }
        } finally {
            setIsSubmitting(false)
        }
    }

    function handleClose() {
        if (isSubmitting) {
            return
        }

        setForm({
            quantity: 1,
            reason: ""
        })

        setErrors({})
        onClose()
    }

    function validateForm(): FormErrors {
        const errors: FormErrors = {}

        if (form.quantity <= 0) {
            errors.quantity = "Quantity must be greater than 0"
        }

        return errors
    }

    return (
        <Dialog
        open={open}
        onClose={handleClose}
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
                error={Boolean(errors.quantity)}
                helperText={errors.quantity ?? "Enter quantity"}
                size="small"
                />

                <TextField
                value={form.reason}
                label="Reason (optional)"
                margin="normal"
                fullWidth
                size="small"
                onChange={(event) =>
                    setForm({
                        ...form,
                        reason: event.target.value,
                    })
                }
                />
            </DialogContent>

            <DialogActions>
                <Button
                onClick={handleClose}
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