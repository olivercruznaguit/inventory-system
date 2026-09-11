import { useEffect, useState } from "react"
import { useAuth } from "../hooks/useAuth"
import { getInventoryDashboard, getRecentStockMovements } from "../services/api"
import type { AlertSeverity, InventoryDashboard, RecentStockMovement } from "../types/dashboard"
import { Alert, Box, CircularProgress, Grid, Typography } from "@mui/material"
import StatCard from "../components/StatCard"
import InventoryAlertCard from "../components/InventoryAlertCard"
import RecentStockMovements from "../components/RecentStockMovements"

export default function Dashboard() {
    const { token } = useAuth()

    const [dashboard, setDashboard] = useState<InventoryDashboard | null>(null)
    const [recentStockMovements, setRecentStockMovements] = useState<RecentStockMovement[]>([])
    
    const [isDashboardLoading, setIsDashboardLoading] = useState(true)
    const [isMovementLoading, setIsMovementLoading] = useState(true)

    const [dashboardError, setDashboardError] = useState(false)
    const [movementError, setMovementError] = useState(false)

    useEffect(() => {
        if (!token) {
            return
        }
        
        const authToken = token

        async function fetchDashboard() {
            try {
                const response = await getInventoryDashboard(authToken)
                setDashboard(response)
            } catch (error) {
                console.error(error)
                setDashboardError(true)
            } finally {
                setIsDashboardLoading(false)
            }
        }

        async function fetchRecentStockMovements() {
            try {
                const response = await getRecentStockMovements(authToken, 10)
                setRecentStockMovements(response.data)
            } catch (error) {
                console.error(error)
                setMovementError(true)
            } finally {
                setIsMovementLoading(false)
            }
        }


        fetchDashboard()
        fetchRecentStockMovements()
    }, [token])
    
    if (isDashboardLoading) {
        return (
            <Box
                sx={{
                    display: "flex",
                    justifyContent: "center",
                    alignItems: "center",
                    minHeight: 300,
                }}
            >
                <CircularProgress />
            </Box>
        )
    }

    if (dashboardError || !dashboard) {
        return (
            <Alert severity="error">
                Failed to load dashboard.
            </Alert>
        )
    }

    const stats = [
        {
            title: "Total Products",
            value: dashboard.totalProducts,
        },
        {
            title: "Active Products",
            value: dashboard.activeProducts,
        },
        {
            title: "Inactive Products",
            value: dashboard.inactiveProducts,
        },
        {
            title: "Quantity on Hand",
            value: dashboard.totalQuantityOnHand,
        },
    ]

    const alerts: {
        title: string
        value: string | number
        description: string
        severity: AlertSeverity
    }[] = [
        {
            title: "Total Inventory Value",
            value: `₱${dashboard.totalInventoryValue.toFixed(2)}`,
            description: "Sum of active products",
            severity: "success"
        },
        {
            title: "Low Stock",
            value: dashboard.lowStockProducts,
            description: "Products approaching minimum stock",
            severity: "warning"
        },
        {
            title: "Out of Stock",
            value: dashboard.outOfStockProducts,
            description: "Products with no stock available",
            severity: "error"
        },
    ]

    return (
        <Box>
            <Typography variant="h4" gutterBottom>
                Dashboard
            </Typography>

            <Typography variant="h6" sx={{ mb: 2 }}>
                Overview
            </Typography>

            <Grid container spacing={2}>
                {stats.map((stat) => (
                    <Grid
                        key={stat.title}
                        size={{ xs: 12, sm: 6, md: 3 }}
                    >
                        <StatCard
                            title={stat.title}
                            value={stat.value}
                        />
                    </Grid>
                ))}
            </Grid>

            <Typography variant="h6" sx={{ mt: 4, mb: 2 }}>
                Inventory Alerts
            </Typography>

            <Grid container spacing={2}>
                {alerts.map(alert => (
                    <Grid 
                    key={alert.title} 
                    size={{ xs: 12, sm: 4 }}>
                        <InventoryAlertCard
                            title={alert.title}
                            count={alert.value}
                            description={alert.description}
                            severity={alert.severity}
                        />
                    </Grid>
                ))}
            </Grid>

            {isMovementLoading && <CircularProgress />}

            {movementError && (
                <Alert severity="error">
                    Failed to load recent stock movements.
                </Alert>
            )}

            {!isMovementLoading && !movementError && (
                <RecentStockMovements stockMovements={recentStockMovements} />
            )}
        </Box>
    )
}