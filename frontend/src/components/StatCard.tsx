import { Card, CardContent, Typography } from "@mui/material"


type StatCardProps = {
    title: string
    value: string | number
}

export default function StatCard({title, value}: StatCardProps){
    return(
        <Card>
            <CardContent>
                <Typography variant="body2">
                    {title}
                </Typography>

                <Typography variant="h4">
                    {value}
                </Typography>
            </CardContent>
        </Card>
    )
}