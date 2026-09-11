import { Card, CardContent, Typography } from "@mui/material"
import type { AlertSeverity } from "../types/dashboard"

type InventoryAlertCardProps = {
    title: string
    count: number | string
    description: string
    severity: AlertSeverity
}

export default function InventoryAlertCard({
    title,
    count,
    description,
    severity,
}: InventoryAlertCardProps) {
    return (
        <Card sx={{ height: "100%" }}>
            <CardContent>
                <Typography
                    variant="body2"
                    sx={{
                        mb: 2,
                        color: "text.secondary",
                    }}
                >
                    {title}
                </Typography>

                <Typography
                    variant="h4"
                    color={severity}
                    sx={{ fontWeight: 700 }}
                >
                    {count}
                </Typography>

                <Typography
                    variant="body2"
                    color={severity}
                >
                    {description}
                </Typography>
            </CardContent>
        </Card>
    )
}