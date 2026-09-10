import { Box, Button, Dialog, DialogActions, DialogContent, DialogTitle, TextField, Typography } from "@mui/material"
import type { Product } from "../types/products"
import { useState } from "react"
import type { InventoryRequest } from "../types/inventory"
import { useAuth } from "../hooks/useAuth"
import { stockOut } from "../services/api"
import { useSnackbar } from "notistack"

type ProductStockOutDialogProp = {
    product: Product | null
    open: boolean
    onClose: () => void
    onSubmit: () => void
}

type FormErrors = {
    quantity?: string
}

export default function ProductStockOutDialog({product, open, onClose, onSubmit}: ProductStockOutDialogProp) {
    const { enqueueSnackbar } = useSnackbar();

    const { token } = useAuth()

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

        const validationErrors = validateForm(product)

        if (Object.keys(validationErrors).length > 0) {
            setErrors(validationErrors)
            return
        }

        const quantity = form.quantity

        try{
            setErrors({})
            setIsSubmitting(true)

            await stockOut(token, product.id, form)

            setForm({
                quantity: 1,
                reason: ""
            })

            onClose()
            onSubmit()

            enqueueSnackbar(
                `Successfully removed ${quantity} ${quantity > 1 ? "units" : "unit" }`,
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

    function validateForm(product: Product): FormErrors {
        const errors: FormErrors = {}

        if (form.quantity <= 0) {
            errors.quantity = "Quantity must be greater than 0"
        } else if (form.quantity > product.quantity) {
            errors.quantity = "Quantity must not be greater than current stock"
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
                Stock Out
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
                    { isSubmitting ? "Stocking out..." : "Stock out" }
                </Button>
            </DialogActions>

        </Dialog>
    )
}