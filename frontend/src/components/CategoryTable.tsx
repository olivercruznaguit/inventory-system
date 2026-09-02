import { Button, LinearProgress, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import type { Category } from "../types/categories";
import TableEmptyState from "./TableEmptyState";

type CategoryTableProps = {
    categories: Category[];
    isFetching: boolean
    onEdit: (category: Category) => void
    onDelete: (category: Category) => void
}

export default function CategoryTable({ categories, isFetching, onEdit, onDelete }: CategoryTableProps) {

    return (
        <TableContainer component={Paper}>
            { isFetching && <LinearProgress aria-label="Fetching…"/>}

            <Table>
                <TableHead>
                    <TableRow>
                        <TableCell>Name</TableCell>
                        <TableCell align="right">Actions</TableCell>
                    </TableRow>
                </TableHead>

                <TableBody>
                    {categories.length === 0 && !isFetching ? (
                        <TableRow>
                            <TableCell colSpan={6} align="center">
                                <TableEmptyState
                                title="No categories yet"
                                description="Add your first category to get started."
                                />
                            </TableCell>
                        </TableRow>
                    ) : (
                    categories.map((category) => (
                        <TableRow  key={category.id}>
                            <TableCell>
                                {category.name}
                            </TableCell>

                            <TableCell align="right">
                                <Button 
                                variant="contained" 
                                onClick={() => onEdit(category)}>
                                    Edit
                                </Button>

                                <Button 
                                variant="contained" 
                                color="error" 
                                sx={{ ml: 1 }}
                                onClick={() => onDelete(category)}>
                                    Delete
                                </Button>
                            </TableCell>
                        </TableRow>
                    )))
                    }

                </TableBody>
            </Table>
        </TableContainer>
    )
}