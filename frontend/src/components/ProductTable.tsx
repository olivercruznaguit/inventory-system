import { Button, Chip, LinearProgress, Menu, MenuItem, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, IconButton } from "@mui/material"
import type { Product } from "../types/products"
import TableEmptyState from "./TableEmptyState"
import { useState } from "react"
import { MoreVert } from "@mui/icons-material"

type ProductTableProps = {
    products: Product[]
    isFetching: boolean
    hasFilters: boolean
    onEdit: (product: Product) => void
    onDelete: (product: Product) => void
    onStockIn: (product: Product) => void
    onStockOut: (product: Product) => void
    onStockHistory: (product: Product) => void
    onResetFilters: () => void
}
export default function ProductTable({ products, isFetching, hasFilters, onEdit, onDelete, onStockIn, onStockOut, onStockHistory, onResetFilters }:ProductTableProps){
    const emptyState = hasFilters
    ? {
        title: "No products found",
        description: "Try adjusting your search or filters.",
      }
    : {
        title: "No products yet",
        description: "Add your first product to get started.",
      }
    
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
    const open = Boolean(anchorEl);

    const handleClick = (event: React.MouseEvent<HTMLElement>, product: Product) => {
        setAnchorEl(event.currentTarget)
        setSelectedProduct(product)
    }

    const onCloseMenu = () => {
        setAnchorEl(null)
        setSelectedProduct(null)
    }

    const handleAction = (action: (p: Product) => void) => {
        if (selectedProduct) {
            action(selectedProduct)
        }
        onCloseMenu()
    };

    return (
        <TableContainer component={Paper}>
            { isFetching && <LinearProgress aria-label="Fetching…"/>}
            <Table size="small">
                <TableHead>
                    <TableRow>
                        <TableCell>Name</TableCell>
                        <TableCell>Category</TableCell>
                        <TableCell>Price</TableCell>
                        <TableCell>Quantity</TableCell>
                        <TableCell>Status</TableCell>
                        <TableCell align="right">Action</TableCell>
                    </TableRow>
                </TableHead>

                <TableBody>
                    {products.length === 0 && !isFetching ? (
                        <TableRow>
                            <TableCell colSpan={6} align="center">
                                <TableEmptyState
                                title={emptyState.title}
                                description={emptyState.description}
                                action={
                                    hasFilters ? (
                                        <Button
                                        variant="outlined" 
                                        onClick={onResetFilters}
                                        >
                                            Reset Filters
                                        </Button>
                                    ) : null
                                }
                                />
                            </TableCell>
                        </TableRow>
                    ) : (
                        products.map((product) => (
                            <TableRow hover key={product.id}>
                                <TableCell>
                                    {product.name}
                                </TableCell>

                                <TableCell>
                                    {product.category?.name ?? "—"}
                                </TableCell>

                                <TableCell>
                                    ₱{product.price.toFixed(2)}
                                </TableCell>

                                <TableCell>
                                    {product.quantity}
                                </TableCell>

                                <TableCell>
                                    { product.status === "ACTIVE" ? 
                                    <Chip color="success" label="Active" /> : 
                                    <Chip color="error" label="Inactive" />  }
                                </TableCell>

                                <TableCell align="right">
                                    <IconButton
                                        size="small"
                                        aria-label={`Actions for ${product.name}`}
                                        id={`button-${product.id}`}
                                        aria-controls={open ? "product-menu" : undefined}
                                        aria-haspopup="true"
                                        aria-expanded={open ? "true" : undefined}
                                        onClick={(e) => handleClick(e, product)}
                                    >
                                        <MoreVert />
                                    </IconButton>
                                </TableCell>
                            </TableRow>
                        ))
                    )}
                </TableBody>
            </Table>

            <Menu
            id="product-menu"
            anchorEl={anchorEl}
            open={open}
            onClose={onCloseMenu}
            >
                <MenuItem onClick={() => handleAction(onEdit)}>Edit</MenuItem>
                <MenuItem onClick={() => handleAction(onStockIn)}>Stock In</MenuItem>
                <MenuItem onClick={() => handleAction(onStockOut)}>Stock Out</MenuItem>
                <MenuItem onClick={() => handleAction(onStockHistory)}>Stock History</MenuItem>
                <MenuItem onClick={() => handleAction(onDelete)}>Delete</MenuItem>
            </Menu>
</TableContainer>
    )
}