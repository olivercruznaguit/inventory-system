import { useEffect, useState } from "react"
import { useAuth } from "../hooks/useAuth"
import { getInventoryDashboard } from "../services/api"
import type { InventoryDashboard } from "../types/dashboard"
import { Box, Grid, Typography } from "@mui/material"
import StatCard from "../components/StatCard"

export default function Dashboard() {
    const { token } = useAuth()

    const [dashboard, setDashboard] = useState<InventoryDashboard | null>(null)
    const [isLoading, setIsLoading] = useState(true)

    useEffect(() => {
        if (!token) {
            return
        }
        const authToken = token
        async function fetchDashboard() {
            try {
                const data = await getInventoryDashboard(authToken)
                setDashboard(data)
            } catch (error) {
                console.error(error)
            } finally {
                setIsLoading(false)
            }
        }

        fetchDashboard()
    }, [token])

    if (isLoading) {
        return <div>Loading...</div>
    }

    if (!dashboard) {
        return <div>Failed to load dashboard</div>
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
        {
            title: "Low Stock",
            value: dashboard.lowStockProducts,
        },
        {
            title: "Out of Stock",
            value: dashboard.outOfStockProducts,
        },
        {
            title: "Total Inventory Value",
            value: `₱${dashboard.totalInventoryValue.toFixed(2)}`,
        },
    ]

    return (
        <Box>
            <Typography variant="h4" gutterBottom>
                Dashboard
            </Typography>

            <Grid container spacing={2}>
                {stats.map((stat) => (
                    <Grid
                        key={stat.title}
                        size={{ xs: 12, sm: 6, md: 4 }}
                    >
                        <StatCard
                            title={stat.title}
                            value={stat.value}
                        />
                    </Grid>
                ))}
            </Grid>
        </Box>
    )
}