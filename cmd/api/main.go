package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/olivercruznaguit/inventory-system/internal/auth"
	"github.com/olivercruznaguit/inventory-system/internal/config"
	"github.com/olivercruznaguit/inventory-system/internal/database"
	"github.com/olivercruznaguit/inventory-system/internal/handler"
	"github.com/olivercruznaguit/inventory-system/internal/middleware"
	"github.com/olivercruznaguit/inventory-system/internal/repository"
	"github.com/olivercruznaguit/inventory-system/internal/service"

	_ "github.com/olivercruznaguit/inventory-system/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Inventory System
// @version         1.0
// @description     REST API for managing products, categories, and inventory stock.
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.email   support@example.com

// @host            localhost:8080
// @BasePath        /
func main() {
	if err := config.LoadEnv(".env"); err != nil {
		log.Fatalf("failed to load env file: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPostgres(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL")

	defer db.Close()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// SWAGGER UI ROUTE
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	productRepository := repository.NewProductRepository(db.DB())
	categoryRepository := repository.NewCategoryRepository(db.DB())
	userRepository := repository.NewUserRepository(db.DB())

	productService := service.NewProductService(productRepository, categoryRepository)
	productHandler := handler.NewProductHandler(productService)

	categoryService := service.NewCategoryService(categoryRepository)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	inventoryService := service.NewInventoryService(db)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)

	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	tokenService := auth.NewTokenService(cfg.Auth.JWTSecret)
	authService := service.NewAuthService(userService, tokenService)
	authHandler := handler.NewAuthHandler(authService)

	authMiddleware := middleware.NewAuthMiddleware(tokenService)

	// PRODUCTS
	products := router.Group("/products")
	products.Use(authMiddleware.Authenticate)
	{
		// list of products
		products.GET("", productHandler.GetProducts)

		// get specific product
		products.GET("/:id", productHandler.GetProductByID)

		// list of stock movements
		products.GET("/:id/stock-movements", inventoryHandler.GetStockMovements)

		// create a new product
		products.POST("", productHandler.CreateProduct)

		// stock-in a specific product
		products.POST("/:id/stock-in", inventoryHandler.StockIn)

		// stock-out a specific product
		products.POST("/:id/stock-out", inventoryHandler.StockOut)

		// update a single product
		products.PUT("/:id", productHandler.UpdateProduct)

		// update the status of a product
		products.PATCH("/:id/status", productHandler.UpdateProductStatus)

		// delete specific product
		products.DELETE("/:id", authMiddleware.RequireAdmin, productHandler.DeleteProduct)
	}

	// CATEGORY
	categories := router.Group("/categories")
	categories.Use(authMiddleware.Authenticate)
	{
		// list of categories
		categories.GET("", categoryHandler.GetCategories)

		// create a category
		categories.POST("", categoryHandler.CreateCategory)

		// get specific category
		categories.GET("/:id", categoryHandler.GetCategoryByID)

		// update a single category
		categories.PUT("/:id", categoryHandler.UpdateCategory)

		// delete a specific category
		categories.DELETE("/:id", authMiddleware.RequireAdmin, categoryHandler.DeleteCategory)
	}

	// INVENTORY
	inventory := router.Group("/inventory")
	inventory.Use(authMiddleware.Authenticate)
	{
		// get the dashboard data
		inventory.GET("/dashboard", inventoryHandler.GetInventoryDashboard)

		// get the recent stock movements
		inventory.GET("/movements", inventoryHandler.GetRecent)
	}

	// USERS
	users := router.Group("/users")
	users.Use(authMiddleware.Authenticate, authMiddleware.RequireAdmin)
	users.POST("", userHandler.CreateUser)

	// AUTH
	router.POST("/auth/login", authHandler.Login)

	fmt.Printf("Server listening on :%s\n", cfg.App.Port)

	if err := router.Run(":" + cfg.App.Port); err != nil {
		log.Fatal(err)
	}
}
