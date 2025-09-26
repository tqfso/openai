package repository

import (
	"context"
	"fmt"
	"openserver/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
)

func Init() error {
	var err error
	cfg := config.GetDatabase()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	pool, err = pgxpool.New(context.Background(), dsn)
	return err
}

func Close() {
	pool.Close()
}

func GetPool() *pgxpool.Pool {
	return pool
}
