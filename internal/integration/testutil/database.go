package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/olivercruznaguit/inventory-system/internal/config"
	"github.com/olivercruznaguit/inventory-system/internal/database"
)

func SetupTestDatabase(t *testing.T) *database.Database {
	t.Helper()

	projectRoot := findProjectRoot(t)
	envPath := filepath.Join(projectRoot, ".env.test")

	if err := config.LoadEnv(envPath); err != nil {
		t.Fatalf("failed to load env file: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	db, err := database.NewPostgres(cfg.DB)
	if err != nil {
		t.Fatalf("failed to create database connection: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func findProjectRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); !os.IsNotExist(err) {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("could not find project root containing go.mod starting from %s", dir)
		}

		dir = parent
	}
}
