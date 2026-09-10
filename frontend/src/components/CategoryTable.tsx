import { Button, IconButton, LinearProgress, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import type { Category } from "../types/categories";
import TableEmptyState from "./TableEmptyState";
import { DeleteOutlineOutlined, EditOutlined } from "@mui/icons-material";

type CategoryTableProps = {
    categories: Category[];
    isFetching: boolean
    hasFilters: boolean
    onEdit: (category: Category) => void
    onDelete: (category: Category) => void
    onResetFilters: () => void
}

export default function CategoryTable({ categories, isFetching, hasFilters, onEdit, onDelete, onResetFilters }: CategoryTableProps) {
    const emptyState = hasFilters
    ? {
        title: "No categories found",
        description: "Try adjusting your search or filters.",
    }
    : {
        title: "No categories yet",
        description: "Add your first category to get started.",
    }


    return (
        <TableContainer component={Paper}>
            { isFetching && <LinearProgress aria-label="Fetching…"/>}

            <Table size="small">
                <TableHead>
                    <TableRow>
                        <TableCell>Name</TableCell>
                        <TableCell>Product Count</TableCell>
                        <TableCell align="right">Actions</TableCell>
                    </TableRow>
                </TableHead>

                <TableBody>
                    {categories.length === 0 && !isFetching ? (
                        <TableRow>
                            <TableCell colSpan={3} align="center">
                                <TableEmptyState
                                    title={emptyState.title}
                                    description={emptyState.description}
                                    action={
                                        hasFilters ? (
                                            <Button
                                                variant="outlined"
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
                    categories.map((category) => (
                        <TableRow  key={category.id}>
                            <TableCell>
                                {category.name}
                            </TableCell>

                            <TableCell>
                                {category.productCount}
                            </TableCell>

                            <TableCell align="right">
                                <IconButton
                                    size="small"
                                    aria-label={`Edit ${category.name}`}
                                    onClick={() => onEdit(category)}
                                >
                                    <EditOutlined />
                                </IconButton>

                                <IconButton
                                    size="small"
                                    color="error"
                                    sx={{ ml: 0.5 }}
                                    aria-label={`Delete ${category.name}`}
                                    disabled={category.productCount > 0}
                                    onClick={() => onDelete(category)}
                                >
                                    <DeleteOutlineOutlined />
                                </IconButton>

                            </TableCell>
                        </TableRow>
                    )))
                    }

                </TableBody>
            </Table>
        </TableContainer>
    )
}