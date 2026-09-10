import { Card, CardContent, Typography } from "@mui/material"

type StatCardProps = {
    title: string
    value: string | number
}

export default function StatCard({ title, value }: StatCardProps) {
    return (
        <Card sx={{ height: "100%" }}>
            <CardContent>
                <Typography
                    variant="body2"
                    
                    sx={{ 
                        mb: 1,
                        color: "text.secondary"
                    }}
                >
                    {title}
                </Typography>

                <Typography
                    variant="h4"
                    color="primary"
                    sx={{ fontWeight: 700 }}
                >
                    {value}
                </Typography>
            </CardContent>
        </Card>
    )
}