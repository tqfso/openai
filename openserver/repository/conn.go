package repository

import (
	"context"
	"fmt"
	"openserver/config"
	"reflect"

	"github.com/jackc/pgx/v5"
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

func IsZeroValue(v any) bool {
	return reflect.ValueOf(v).IsZero()
}

func WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	conn, err := GetPool().Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
