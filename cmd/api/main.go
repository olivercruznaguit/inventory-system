package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/olivercruznaguit/inventory-system/internal/config"
	"github.com/olivercruznaguit/inventory-system/internal/database"
	"github.com/olivercruznaguit/inventory-system/internal/handler"
	"github.com/olivercruznaguit/inventory-system/internal/repository"
	"github.com/olivercruznaguit/inventory-system/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
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

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	router.GET("/products", productHandler.GetProducts)
	router.GET("/products/:id", productHandler.GetProductByID)
	router.POST("/products", productHandler.CreateProduct)
	router.PUT("/products/:id", productHandler.UpdateProduct)
	router.PATCH("/products/:id/status", productHandler.UpdateProductStatus)
	router.DELETE("/products/:id", productHandler.DeleteProduct)

	fmt.Printf("Server listening on :%s\n", cfg.App.Port)

	if err := router.Run(":" + cfg.App.Port); err != nil {
		log.Fatal(err)
	}
}
