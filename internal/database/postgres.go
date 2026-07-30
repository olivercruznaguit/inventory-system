package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/olivercruznaguit/inventory-system/internal/config"
)

type Database struct {
	pool *pgxpool.Pool
}

func NewPostgres(cfg config.DatabaseConfig) (*Database, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(context.Background()); err != nil {
		db.Close()
		return nil, err
	}

	return &Database{
		pool: db,
	}, nil
}

func (db *Database) DB() *pgxpool.Pool {
	return db.pool
}

func (db *Database) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return db.pool.Begin(ctx)
}

func (db *Database) Close() {
	db.pool.Close()
}
