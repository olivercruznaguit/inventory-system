import { Alert, Box, Button, Chip, CircularProgress, Dialog, DialogActions, DialogContent, DialogTitle, Table, TableBody, TableCell, TableHead, TableRow, Typography } from "@mui/material"
import { useAuth } from "../hooks/useAuth"
import { useCallback, useEffect, useState } from "react"
import type { StockMovement } from "../types/inventory"
import { stockMovement } from "../services/api"
import type { Product } from "../types/products"

type ProductStockHistoryDialogProps = {
    open: boolean
    product: Product | null
    onClose: () => void
}

export default function ProductStockHistoryDialog({ open, product, onClose }: ProductStockHistoryDialogProps) {
    const { token } = useAuth()

    const [isFetching, setIsFetching] = useState(true)

    const [error, setError] = useState<string | null>(null)

    const [productStockMovements, setProductStockMovements] = useState<StockMovement[] | null>(null)

    const formatDate = (date: string) => {
        if (!date) return '—';
        return new Date(date).toLocaleString(undefined, {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
            hour12: true
        });
    };

    const fetchStockMovements = useCallback(async ()=>{
         if(!open || !token || !product) {
            return
        }

        try {
            setIsFetching(true)
            setError(null)
            setProductStockMovements(null)

            const response = await stockMovement(token, product.id)

            setProductStockMovements(response.data)
        } catch (error) {
            console.error(error)
            setError("Failed to fetch product stock history")
        } finally {
            setIsFetching(false)
        }
    }, [token, product, open])

    useEffect(() => {
        fetchStockMovements()
    }, [fetchStockMovements])

    return (
        <Dialog 
        open={open} 
        onClose={onClose}
        fullWidth
        maxWidth="lg"
        >
            <DialogTitle>
                Stock History
            </DialogTitle>

            <DialogContent>
                { isFetching && 
                    <CircularProgress/>
                }

                { error && 
                    <Alert severity="error">
                        { error }
                    </Alert>
                }

                { productStockMovements && 
                    <Box>

                        <Box 
                        sx={{ 
                            display: "flex",
                            flexDirection: "column",
                            gap: 1
                        }}
                        >
                            <Box
                            sx={{ 
                                display: "flex",
                                gap: 1
                            }}
                            >
                                <Typography variant="subtitle1">
                                    Product:
                                </Typography>

                                <Typography variant="subtitle1" sx={{fontWeight: "bold"}}>
                                    {product?.name}
                                </Typography>
                            </Box>

                            <Box
                            sx={{ 
                                display: "flex",
                                gap: 1
                            }}
                            >
                                <Typography variant="subtitle1">
                                    Current Stock:
                                </Typography>

                                <Typography variant="subtitle1" sx={{fontWeight: "bold"}}>
                                    {product?.quantity}
                                </Typography>
                            </Box>
                        </Box>

                        <Table>
                            <TableHead>
                                <TableRow>
                                    <TableCell>
                                        Date    
                                    </TableCell>
                                    
                                    <TableCell>
                                        Type
                                    </TableCell>
                                    
                                    <TableCell>
                                        Change
                                    </TableCell>
                                    
                                    <TableCell>
                                        Stock After
                                    </TableCell>

                                    <TableCell>
                                        Reason
                                    </TableCell>
                                </TableRow>
                            </TableHead>

                            <TableBody>
                                { productStockMovements.map((stockMovement) => (
                                    <TableRow key={stockMovement.id}>
                                        <TableCell>
                                            {formatDate(stockMovement.createdAt)}    
                                        </TableCell>
                                        
                                        <TableCell>
                                            { stockMovement.type === "IN" ? 
                                            <Chip color="success" label="IN" /> : 
                                            <Chip color="error" label="OUT" />  }
                                        </TableCell>
                                        
                                        <TableCell>
                                            <Typography variant="body1" color={stockMovement.type === "IN" ? "success" : "error"}>
                                                {stockMovement.type === "IN" ? 
                                                "+" : "-"
                                                }

                                                {stockMovement.quantity}
                                            </Typography>
                                        </TableCell>
                                        
                                        <TableCell>
                                            {stockMovement.remainingQuantity}
                                        </TableCell>

                                        <TableCell>
                                            {stockMovement.reason || "—"}
                                        </TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </Box>

                }
            </DialogContent>

            <DialogActions>
                <Button
                onClick={onClose}
                >
                    Close
                </Button>
            </DialogActions>
        </Dialog>
    )
}