import { Box, Typography } from "@mui/material"
import type { ReactNode } from "react"

type TableEmptyStateProps = {
    title: string
    description: string
    action?: ReactNode
}

export default function TableEmptyState({ title, description, action }: TableEmptyStateProps) {
    return (
        <Box sx={{ p: 5, textAlign: "center" }}>
            <Typography variant="h6" sx={{ mb: 1 }}>
                {title}
            </Typography>
            <Typography variant="body2" sx={{ mb: 2 }}>
                {description}
            </Typography>
            {action}
        </Box>
    )
}