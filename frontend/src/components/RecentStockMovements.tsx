import { Button, Chip, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Typography } from "@mui/material"
import type { RecentStockMovement } from "../types/dashboard"
import TableEmptyState from "./TableEmptyState"
import { Link } from "react-router-dom"

type RecentStockMovementsProps = {
    stockMovements: RecentStockMovement[]
}

export default function RecentStockMovements({ stockMovements }: RecentStockMovementsProps){
    const formatDate = (date: string) => {
        if (!date) return '—'

        return new Date(date).toLocaleString(undefined, {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
            hour12: true
        })
    }
    return (
        <Paper sx={{ 
            mt: 4,
            p: 2,
            borderRadius: "10px",
        }}>
            <Typography variant="h6" gutterBottom>
                Recent Stock Movements
            </Typography>

            <TableContainer sx={{ overflowX: 'auto', maxWidth: '100%' }}>
                <Table sx={{ minWidth: 800 }}>
                    <TableHead>
                        <TableRow>
                            <TableCell>
                                Product
                            </TableCell>

                            <TableCell>
                                Reason
                            </TableCell>

                            <TableCell>
                                Type
                            </TableCell>

                            <TableCell align="right">
                                Quantity
                            </TableCell>

                            <TableCell align="right">
                                Stock After
                            </TableCell>

                            <TableCell>
                                Date
                            </TableCell>
                        </TableRow>
                    </TableHead>

                    <TableBody>
                        {stockMovements.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={6} align="center">
                                    <TableEmptyState
                                        title="No Stock Movements Yet"
                                        description="Stock movements will appear here when inventory changes"
                                        action={
                                            <Button component={Link} to="/products" variant="outlined">
                                                Go to Products
                                            </Button>
                                        }
                                        />
                                </TableCell>
                            </TableRow>
                        ):
                        (stockMovements.map((stockMovement) => (
                            <TableRow key={stockMovement.id}>
                                <TableCell>
                                    { stockMovement.productName }
                                </TableCell>
                                
                                <TableCell>
                                    { stockMovement.reason || "—" }
                                </TableCell>

                                <TableCell>
                                    { stockMovement.type === "IN" ? 
                                    <Chip size="small" color="success" label="IN" /> : 
                                    <Chip size="small" color="error" label="OUT" />  }
                                </TableCell>

                                <TableCell align="right">
                                    <Typography variant="body2" color={stockMovement.type === "IN" ? "success" : "error"}>
                                        {stockMovement.type === "IN" ? 
                                        "+" : "-"
                                        }

                                        { stockMovement.quantity }
                                    </Typography>
                                </TableCell>

                                <TableCell align="right">
                                    { stockMovement.remainingQuantity }
                                </TableCell>
                                
                                <TableCell>
                                    { formatDate(stockMovement.createdAt) }
                                </TableCell>
                            </TableRow>
                        )))}
                    </TableBody>
                </Table>
            </TableContainer>
        </Paper>
    )
}