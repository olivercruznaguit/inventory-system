import { Button, Dialog, DialogActions, DialogContent, DialogTitle, FormControl, InputLabel, MenuItem, Select, TextField } from "@mui/material"
import type { Category } from "../types/categories"
import { useAuth } from "../hooks/useAuth"
import type { CreateProductRequest } from "../types/products"
import { useState } from "react"
import { createProduct } from "../services/api"
import { useSnackbar } from "notistack"

type CreateProductDialogProps = {
    open: boolean
    categories: Category[]
    onClose: () => void
    onCreated: () => void
}

type FormErrors = {
    name?: string
    price?: string
    minimumStock?: string
}

export default function CreateProductDialog({open, categories, onClose, onCreated}:CreateProductDialogProps){
    const { token } = useAuth()

    const { enqueueSnackbar } = useSnackbar()

    const [form, setForm] = useState<CreateProductRequest>({
        name: "",
        price: 0,
        minimumStock: 0,
    })

    const [isCreating, setIsCreating] = useState(false)

    const [errors, setErrors] = useState<FormErrors>({})

    async function handleSubmit() {
        if (!token) {
            return
        }

        const validationErrors = validateForm()

        if (Object.keys(validationErrors).length > 0) {
            setErrors(validationErrors)
            return
        }

        try {
            setIsCreating(true)
            setErrors({})

            const request = {
                ...form,
                name: form.name.trim(),
            }

            await createProduct(token, request)

            setForm({
                name: "",
                price: 0,
                minimumStock: 0,
            })

            onClose()
            onCreated()

            enqueueSnackbar("Product added", { variant: "success" })
        } catch (error) {
            console.error(error)
            enqueueSnackbar("Failed to add product", { variant: "error" })
        } finally {
            setIsCreating(false)
        }
    }

    function handleClose() {
        if(isCreating) {
            return
        }

        setForm({
            name: "",
            price: 0,
            minimumStock: 0,
        })

        setErrors({})
        onClose()
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

    return (
        <Dialog
            open={open}
            onClose={handleClose}
            fullWidth
            maxWidth="sm"
        >
            <DialogTitle>
                Create Product
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
                    size="small"
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
                    error={Boolean(errors.price)}
                    helperText={errors.price}
                />

                <FormControl fullWidth size="small" margin="normal">
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

                <TextField
                    fullWidth
                    size="small"
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
                    disabled={isCreating}
                    onClick={handleSubmit}
                >
                    {isCreating ? "Creating..." : "Create"}
                </Button>
            </DialogActions>
        </Dialog>
    )
}