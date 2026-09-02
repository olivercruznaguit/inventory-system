import { Box, Pagination, Typography } from "@mui/material"

type CustomPaginationProps = {
    page: number
    totalPages: number
    totalItems: number
    currentItemCount: number
    onPageChange: (page: number) => void
}


export default function CustomPagination({page, totalPages, totalItems, currentItemCount, onPageChange}: CustomPaginationProps){
    
    if (totalItems <= 0) {
        return null
    }
    
    return (
        <Box sx={{
            display: "flex",
            justifyContent: "space-between",
            mt: 3
        }}>

            <Typography
                variant="body2"
                sx={{ mb: 2 }}
            >
                Showing {currentItemCount} of{" "}
                {totalItems} items
            </Typography>

            {totalPages > 1 && (
                <Pagination
                    count={totalPages}
                    page={page}
                    onChange={(_, value) => onPageChange(value)}
                />
            )}
        </Box>
    )
}