import { Button, Dialog, DialogActions, DialogContent, DialogTitle, FormControl, InputLabel, MenuItem, Select, TextField, Typography } from "@mui/material"
import type { Category } from "../types/categories"
import { useAuth } from "../hooks/useAuth"
import type { CreateProductRequest } from "../types/products"
import { useState } from "react"
import { createProduct } from "../services/api"

type CreateProductDialogProps = {
    open: boolean
    categories: Category[]
    onClose: () => void
    onCreated: () => void
}

export default function CreateProductDiallog({open, categories, onClose, onCreated}:CreateProductDialogProps){
    const { token } = useAuth()

    const [form, setForm] = useState<CreateProductRequest>({
        name: "",
        price: 0,
        minimumStock: 0,
    })

    const [isCreating, setIsCreating] = useState(false)

    const [error, setError] = useState<string | null>(null)

    async function handleSubmit() {
        if (!token) {
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
            setIsCreating(true)
            setError(null)

            await createProduct(token, form)

            setForm({
                name: "",
                price: 0,
                minimumStock: 0,
            })

            onClose()
            onCreated()
        } catch (error) {
            console.error(error)
            setError("Failed to create product")
        } finally {
            setIsCreating(false)
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
                Create Product
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
                    disabled={isCreating}
                    onClick={handleSubmit}
                >
                    {isCreating ? "Creating..." : "Create"}
                </Button>
            </DialogActions>
        </Dialog>
    )
}