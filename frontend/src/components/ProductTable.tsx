import { Button, LinearProgress, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material"
import type { Product } from "../types/products"
import TableEmptyState from "./TableEmptyState"

type ProductTableProps = {
    products: Product[]
    isFetching: boolean
    hasFilters: boolean
    onEdit: (product: Product) => void
    onDelete: (product: Product) => void
    onResetFilters: () => void
}
export default function ProductTable({ products, isFetching, hasFilters, onEdit, onDelete, onResetFilters }:ProductTableProps){
    const emptyState = hasFilters
    ? {
        title: "No products found",
        description: "Try adjusting your search or filters.",
      }
    : {
        title: "No products yet",
        description: "Add your first product to get started.",
      }
    
    return (
        <TableContainer component={Paper}>
            { isFetching && <LinearProgress aria-label="Fetching…"/>}
            <Table>
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
                                        variant="contained" 
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
                            <TableRow key={product.id}>
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
                                    {product.status}
                                </TableCell>

                                <TableCell align="right">
                                    <Button
                                    variant="contained"
                                    onClick={() => onEdit(product)}>
                                        Edit
                                    </Button>
                                    <Button
                                    variant="contained"
                                    color="error"
                                    sx={{
                                        ml: 1
                                    }}
                                    onClick={() => onDelete(product)}>
                                        Delete
                                    </Button>
                                </TableCell>
                            </TableRow>
                        ))
                    )}
                </TableBody>
            </Table>
        </TableContainer>
    )
}