package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/olivercruznaguit/inventory-system/internal/auth"
	"github.com/olivercruznaguit/inventory-system/internal/config"
	"github.com/olivercruznaguit/inventory-system/internal/database"
	"github.com/olivercruznaguit/inventory-system/internal/handler"
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

	// PRODUCTS
	router.GET("/products", productHandler.GetProducts)
	router.GET("/products/:id", productHandler.GetProductByID)
	router.POST("/products", productHandler.CreateProduct)
	router.PUT("/products/:id", productHandler.UpdateProduct)
	router.PATCH("/products/:id/status", productHandler.UpdateProductStatus)
	router.DELETE("/products/:id", productHandler.DeleteProduct)

	// CATEGORY
	router.POST("/categories", categoryHandler.CreateCategory)
	router.GET("/categories", categoryHandler.GetCategories)
	router.GET("/categories/:id", categoryHandler.GetCategoryByID)
	router.PUT("/categories/:id", categoryHandler.UpdateCategory)
	router.DELETE("/categories/:id", categoryHandler.DeleteCategory)

	// INVENTORY
	router.GET("/products/:id/stock-movements", inventoryHandler.GetStockMovements)
	router.GET("/inventory/dashboard", inventoryHandler.GetInventoryDashboard)
	router.POST("/products/:id/stock-in", inventoryHandler.StockIn)
	router.POST("/products/:id/stock-out", inventoryHandler.StockOut)

	// AUTH
	router.POST("/auth/login", authHandler.Login)
	router.POST("/auth/register", userHandler.Register)

	fmt.Printf("Server listening on :%s\n", cfg.App.Port)

	if err := router.Run(":" + cfg.App.Port); err != nil {
		log.Fatal(err)
	}
}
