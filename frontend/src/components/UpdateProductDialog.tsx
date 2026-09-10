import { Button, Dialog, DialogActions, DialogContent, DialogTitle, FormControl, InputLabel, MenuItem, Select, TextField } from "@mui/material"
import type { Category } from "../types/categories"
import type { Product, UpdateProductRequest } from "../types/products"
import { useState } from "react"
import { useAuth } from "../hooks/useAuth"
import { updateProduct } from "../services/api"
import { useSnackbar } from "notistack"

type UpdateProductDialogProps = {
    open: boolean
    product: Product | null
    categories: Category[]
    onClose: () => void
    onUpdated: () => void
}

type FormErrors = {
    name?: string
    price?: string
    minimumStock?: string
}

export default function UpdateProductDialog({open, product, categories, onClose, onUpdated}: UpdateProductDialogProps) {
    const { token } = useAuth()

    const { enqueueSnackbar } = useSnackbar()

    const [form, setForm] = useState<UpdateProductRequest>({
        name: product?.name || "",
        price: product?.price ?? 0,
        status: product?.status ?? "ACTIVE",
        categoryId: product?.category?.id,
        minimumStock: product?.minimumStock ?? 0
    })

    const [errors, setErrors] = useState<FormErrors>({})
   
    const [isUpdating, setIsUpdating] = useState(false)

    async function handleSubmit() {
        if (!token || !product) {
            return
        }

        const validationErrors = validateForm()

        if (Object.keys(validationErrors).length > 0) {
            setErrors(validationErrors)
            return
        }

        try {
            setIsUpdating(true)
            setErrors({})

            const request = {
                ...form,
                name: form.name.trim(),
            }

            await updateProduct(token, product.id, request)

            onClose()
            onUpdated()

            enqueueSnackbar("Product updated", { variant: "success" })
        } catch (error) {
            console.error(error)
            enqueueSnackbar("Failed to update product", { variant: "error" })
        } finally {
            setIsUpdating(false)
        }
    }

    function validateForm(): FormErrors {
        const errors: FormErrors = {}

        if (!form.name.trim()) {
            errors.name = "Product name is required"
        }

        if (form.price <= 0) {
            errors.price = "Price must be greater than 0"
        }

        if (form.minimumStock < 0) {
            errors.minimumStock = "Minimum stock cannot be negative"
        }

        return errors
    }

    function handleClose() {
        if (isUpdating) {
            return
        }

        setErrors({})
        onClose()
    }

    return (
        <Dialog 
        open={open} 
        onClose={handleClose} 
        fullWidth 
        maxWidth="sm">
            <DialogTitle>
                Update Product
            </DialogTitle>

            <DialogContent>
                <TextField
                    required
                    fullWidth
                    size="small"
                    label="Name"
                    margin="normal"
                    value={form.name}
                    onChange={(event) =>
                        setForm({
                            ...form,
                            name: event.target.value,
                        })
                    }
                    error={Boolean(errors.name)}
                    helperText={errors.name}
                />

                <TextField
                    fullWidth
                    label="Price"
                    size="small"
                    type="number"
                    margin="normal"
                    value={form.price}
                    onChange={(event) =>
                        setForm({
                            ...form,
                            price: Number(event.target.value),
                        })
                    }
                    error={Boolean(errors.price)}
                    helperText={errors.price}
                />

                <FormControl 
                fullWidth 
                size="small"
                margin="normal">
                    <InputLabel>Category</InputLabel>

                    <Select
                        value={form.categoryId?.toString() ?? ""}
                        label="Category"
                        onChange={(event) => {
                            const value = event.target.value

                            setForm({
                                ...form,
                                categoryId:
                                    value === ""
                                        ? undefined
                                        : Number(value),
                            })
                        }}
                    >
                        <MenuItem value="">
                            No Category
                        </MenuItem>

                        {categories.map((category) => (
                            <MenuItem
                                key={category.id}
                                value={category.id}
                            >
                                {category.name}
                            </MenuItem>
                        ))}
                    </Select>
                </FormControl>

                <FormControl 
                fullWidth 
                size="small"
                margin="normal">
                    <InputLabel>Status</InputLabel>

                    <Select
                        value={form.status ?? "ACTIVE"}
                        label="Status"
                        onChange={(event) =>
                            setForm({
                                ...form,
                                status: event.target.value,
                            })
                        }
                    >
                        <MenuItem
                            key="ACTIVE"
                            value="ACTIVE"
                        >
                            ACTIVE
                        </MenuItem>
                        <MenuItem
                            key="INACTIVE"
                            value="INACTIVE"
                        >
                            INACTIVE
                        </MenuItem>
                    </Select>
                </FormControl>

                <TextField
                    fullWidth
                    label="Minimum Stock"
                    type="number"
                    size="small"
                    margin="normal"
                    value={form.minimumStock}
                    onChange={(event) =>
                        setForm({
                            ...form,
                            minimumStock: Number(event.target.value),
                        })
                    }
                    error={Boolean(errors.minimumStock)}
                    helperText={errors.minimumStock}
                />

            </DialogContent>

             <DialogActions>
                <Button onClick={handleClose}>
                    Cancel
                </Button>

                <Button
                    variant="contained"
                    disabled={isUpdating}
                    onClick={handleSubmit}
                >
                    {isUpdating ? "Updating..." : "Update"}
                </Button>
            </DialogActions>

        </Dialog>
    )
}