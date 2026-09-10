import { useCallback, useEffect, useState } from "react"
import { useAuth } from "../hooks/useAuth"
import { getCategories, getProducts } from "../services/api"
import type { Product, ProductPagination, ProductQueryParams, ProductSort, ProductStatus } from "../types/products"
import { Alert, Box, Button, CircularProgress, Typography } from "@mui/material"
import type { Category, CategoryQueryParams } from "../types/categories"
import CreateProductDialog from "../components/CreateProductDialog"
import UpdateProductDialog from "../components/UpdateProductDialog"
import DeleteProductDialog from "../components/DeleteProductDialog"
import { productSortParams } from "../constants/products"
import ProductFilters from "../components/ProductFilters"
import ProductTable from "../components/ProductTable"
import CustomPagination from "../components/CustomPagination"
import ProductStockInDialog from "../components/ProductStockInDialog"
import ProductStockOutDialog from "../components/ProductStockOutDialog"
import ProductStockHistoryDialog from "../components/ProductStockHistoryDialog"

export default function Products() {
    const { token } = useAuth()

    const [products, setProducts] = useState<Product[]>([])
    
    const [isLoading, setIsLoading] = useState<boolean>(true)

    const [isFetching, setIsFetching] = useState(false)

    const [pagination, setPagination] = useState<ProductPagination | null>(null)

    const [error, setError] = useState<string | null>(null)

    const [page, setPage] = useState(1)

    const [search, setSearch] = useState("")

    const [status, setStatus] = useState<ProductStatus>("ACTIVE")

    const [categoryFilter, setCategoryFilter] = useState("")

    const [debouncedSearch, setDebouncedSearch] = useState("")

    const [productSort, setProductSort] = useState<ProductSort>("")

    const [openCreateDialog, setOpenCreateDialog] = useState(false)
    
    const [openUpdateDialog, setOpenUpdateDialog] = useState(false)

    const [openDeleteDialog, setOpenDeleteDialog] = useState(false)

    const [openStockInDialog, setOpenStockInDialog] = useState(false)

    const [openStockOutDialog, setOpenStockOutDialog] = useState(false)

    const [openStockHistoryDialog, setOpenStockHistoryDialog] = useState(false)

    const [categories, setCategories] = useState<Category[]>([])
    
    const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)

    const pageSize = 10
    


    const handleOpenUpdateProduct = (product: Product) => {
        setSelectedProduct(product)
        setOpenUpdateDialog(true)
    }
    
    const handleOpenDeleteProduct = (product: Product) => {
        setSelectedProduct(product)
        setOpenDeleteDialog(true)
    }

    const handleOpenStockInProduct = (product: Product) => {
        setSelectedProduct(product)
        setOpenStockInDialog(true)
    }

    const handleOpenStockOutProduct = (product: Product) => {
        setSelectedProduct(product)
        setOpenStockOutDialog(true)
    }

    const handleOpenStockHistoryProduct = (product: Product) => {
        setSelectedProduct(product)
        setOpenStockHistoryDialog(true)
    }

    const handleProductDeleted = () => {
        if(pagination) {
            const newTotalItems = pagination?.totalItems - 1
            const newTotalPages = Math.max(1, Math.ceil(newTotalItems / pageSize))
            
            if (page > newTotalPages) {
                setPage(newTotalPages)
            } else {
                fetchProducts()
            }
        } else {
            setPage(1)
        }
    }

    const handleSearchChange = (searchStr: string) => {
        setSearch(searchStr)
        setPage(1)
    }

    const handleStatusChange = (sts: ProductStatus) => {
        setStatus(sts)
        setPage(1)
    }

    const handleCategoryChange = (categoryID: string) => {
        setCategoryFilter(categoryID)
        setPage(1)
    }

    const handleSortChange = (sort: ProductSort) => {
        setProductSort(sort)
        setPage(1)
    }

    const handleResetFilters = () => {
        setSearch("")
        setDebouncedSearch("")
        setStatus("ACTIVE")
        setCategoryFilter("")
        setProductSort("")
        setPage(1)
    }

    const hasFilters =
    search.trim() !== "" ||
    status !== "ACTIVE" ||
    categoryFilter !== "" ||
    productSort !== ""


    const fetchProducts = useCallback(async (signal?: AbortSignal) => {
        if (!token) {
            return
        }

        try {
            setIsFetching(true)
            setError(null)

            const {sortBy, sortOrder} = productSortParams[productSort]

            const params: ProductQueryParams = {
                page,
                pageSize,
                status,
                search: debouncedSearch,
                categoryId: categoryFilter,
                sortBy,
                sortOrder
            }

            const response = await getProducts(
                token,
                params,
                signal
            )

            setProducts(response.data)
            setPagination(response.pagination)
        } catch (error) {

            if (error instanceof Error && error.name === 'AbortError') {
                console.log('Fetch successfully aborted')
                return; 
            }

            console.error(error)
            setError("Failed to load products")
        } finally {
            if (!signal?.aborted) {
                setIsLoading(false)
                setIsFetching(false)
            }
        }
    }, [token, page, status, debouncedSearch, categoryFilter, productSort])

    useEffect(() => {
        const controller = new AbortController()

        fetchProducts(controller.signal)

        return () => {
            controller.abort()
        };
    }, [fetchProducts])

    useEffect(() => {
        if(!token){
            return
        }

        const authToken = token

        async function fetchCategories() {
            try {
                const params: CategoryQueryParams = {
                    page: 1,
                    pageSize: 100, // Adjust as needed
                    search: "",
                    sortBy: "",
                    sortOrder: "" 
                }

                const response = await getCategories(authToken, params)
                setCategories(response.data)
            } catch (error) {
                console.error(error)
            }

        }

        fetchCategories()
    },[token])

    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedSearch(search)
        }, 500)

        return () => {
            clearTimeout(timer)
        }
    }, [search])

    if (isLoading) {
        return (
            <Box
                sx={{
                    display: "flex",
                    justifyContent: "center",
                    alignItems: "center",
                    minHeight: "300px",
                }}
            >
                <CircularProgress />
            </Box>
        )
    }

    if (error) {
        return <Alert severity="error">{error}</Alert>
    }

    return (
        <Box>
            <Box
                sx={{
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "center",
                    mb: 2,
                }}
            >
                <Box
                    sx={{
                        display: "flex",
                        alignItems: "center",
                        gap: 1,
                    }}
                >
                    <Typography variant="h4">
                        Products
                    </Typography>
                </Box>
                

                <Button
                    variant="contained"
                    onClick={() => {
                        setOpenCreateDialog(true)
                    }}
                >
                    Add Product
                </Button>
            </Box>

            <ProductFilters
                search={search}
                status={status}
                categoryFilter={categoryFilter}
                productSort={productSort}
                categories={categories}
                onSearchChange={handleSearchChange}
                onStatusChange={handleStatusChange}
                onCategoryChange={handleCategoryChange}
                onSortChange={handleSortChange}
                onResetFilters={handleResetFilters}
            />


            <ProductTable 
            isFetching={isFetching}
            products={products}
            onEdit={handleOpenUpdateProduct}
            onDelete={handleOpenDeleteProduct}
            onStockIn={handleOpenStockInProduct}
            onStockOut={handleOpenStockOutProduct}
            onStockHistory={handleOpenStockHistoryProduct}
            hasFilters={hasFilters}
            onResetFilters={handleResetFilters}
            />

            {pagination &&     
                <CustomPagination 
                page={page} 
                totalPages={pagination.totalPages} 
                totalItems={pagination.totalItems} 
                currentItemCount={products.length} 
                onPageChange={setPage}                
                />
            }

            <CreateProductDialog
                open={openCreateDialog}
                categories={categories}
                onClose={() => setOpenCreateDialog(false)}
                onCreated={() => {
                    fetchProducts()
                    setSelectedProduct(null)
                }}
            />

            <UpdateProductDialog
                key={selectedProduct?.id ?? 'update-product'}
                open={openUpdateDialog}
                categories={categories}
                product={selectedProduct}
                onClose={() => setOpenUpdateDialog(false)}
                onUpdated={() => {
                    fetchProducts()
                    setSelectedProduct(null)
                }}
            />

            <DeleteProductDialog
                open={openDeleteDialog}
                product={selectedProduct}
                onClose={() => setOpenDeleteDialog(false)}
                onDeleted={handleProductDeleted}
            />

            <ProductStockInDialog 
            key={`stock-in-${selectedProduct?.id ?? "closed"}`}
            open={openStockInDialog}
            product={selectedProduct}
            onClose={() => setOpenStockInDialog(false)}
            onSubmit={()=> {
                fetchProducts()
                setSelectedProduct(null)
            }}
            />

            <ProductStockOutDialog 
            key={`stock-out-${selectedProduct?.id ?? "closed"}`}
            open={openStockOutDialog}
            product={selectedProduct}
            onClose={() => setOpenStockOutDialog(false)}
            onSubmit={()=> {
                fetchProducts()
                setSelectedProduct(null)
            }}
            />

            <ProductStockHistoryDialog 
            open={openStockHistoryDialog}
            product={selectedProduct}
            onClose={() => {
                setOpenStockHistoryDialog(false)
            }}
            />
        </Box>
    )
}