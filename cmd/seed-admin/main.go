package main

import (
	"context"
	"log"

	"github.com/olivercruznaguit/inventory-system/internal/config"
	"github.com/olivercruznaguit/inventory-system/internal/database"
	"github.com/olivercruznaguit/inventory-system/internal/repository"
	"github.com/olivercruznaguit/inventory-system/internal/service"
)

func main() {
	if err := config.LoadEnv(".env"); err != nil {
		log.Fatalf("failed to load env file: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	adminEmail := cfg.Bootstrap.AdminEmail
	adminPassword := cfg.Bootstrap.AdminPassword
	if adminEmail == "" {
		log.Fatal("ADMIN_EMAIL is required")
	}

	if adminPassword == "" {
		log.Fatal("ADMIN_PASSWORD is required")
	}

	db, err := database.NewPostgres(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL")

	defer db.Close()

	userRepository := repository.NewUserRepository(db.DB())
	userService := service.NewUserService(userRepository)

	user, err := userService.CreateAdmin(context.Background(), adminEmail, adminPassword)

	if err != nil {
		log.Fatalf("failed to create admin: %v", err)
	}

	log.Printf("created admin user: %s", user.Email)
}
