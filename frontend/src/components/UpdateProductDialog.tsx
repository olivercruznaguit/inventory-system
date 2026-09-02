import { Button, Dialog, DialogActions, DialogContent, DialogTitle, FormControl, InputLabel, MenuItem, Select, TextField, Typography } from "@mui/material"
import type { Category } from "../types/categories"
import type { Product, UpdateProductRequest } from "../types/products"
import { useState } from "react"
import { useAuth } from "../hooks/useAuth"
import { updateProduct } from "../services/api"

type UpdateProductDialogProps = {
    open: boolean
    product: Product | null
    categories: Category[]
    onClose: () => void
    onUpdated: () => void
}

export default function UpdateProductDialog({open, product, categories, onClose, onUpdated}: UpdateProductDialogProps) {
    const { token } = useAuth()

    const [form, setForm] = useState<UpdateProductRequest | null>({
        name: product?.name || "",
        price: product?.price ?? 0,
        status: product?.status ?? "ACTIVE",
        categoryId: product?.category?.id,
        minimumStock: product?.minimumStock ?? 0
    })

    const [error, setError] = useState<string | null>(null)
   
    const [isUpdating, setIsUpdating] = useState<boolean>(false)

    async function handleSubmit() {
        if (!token || !form || !product) {
            return
        }

        if (!form.name.trim()) {
            setError("Product name is required")
            return
        }

        if (form.price <= 0) {
            setError("Price must be greater than 0")
            return
        }

        try {
            setIsUpdating(true)
            setError(null)

            await updateProduct(token, product.id, form)

            setForm(null)

            onClose()
            onUpdated()
        } catch (error) {
            console.error(error)
            setError("Failed to update product")
        } finally {
            setIsUpdating(false)
        }
    }

    if (!form) return null;

    return (
        <Dialog 
        open={open} 
        onClose={onClose} 
        fullWidth 
        maxWidth="sm">
            <DialogTitle>
                Update Product
            </DialogTitle>

            <DialogContent>
                <TextField
                    fullWidth
                    label="Name"
                    margin="normal"
                    value={form.name}
                    onChange={(event) =>
                        setForm({
                            ...form,
                            name: event.target.value,
                        })
                    }
                />

                <TextField
                    fullWidth
                    label="Price"
                    type="number"
                    margin="normal"
                    value={form.price}
                    onChange={(event) =>
                        setForm({
                            ...form,
                            price: Number(event.target.value),
                        })
                    }
                />

                <FormControl fullWidth margin="normal">
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

                <FormControl fullWidth margin="normal">
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
                    margin="normal"
                    value={form.minimumStock}
                    onChange={(event) =>
                        setForm({
                            ...form,
                            minimumStock: Number(event.target.value),
                        })
                    }
                />

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
                <Button onClick={onClose}>
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