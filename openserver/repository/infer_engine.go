package repository

import (
	"context"
	"errors"
	"openserver/model"

	"github.com/jackc/pgx/v5"
)

type InferEngineRepo struct{}

func InferEngine() *InferEngineRepo {
	return &InferEngineRepo{}
}

func (r *InferEngineRepo) GetByName(ctx context.Context, name string) (*model.InferEngine, error) {
	conn, err := GetPool().Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, `
		SELECT name, framework, image, status, created_at, updated_at 
		FROM infer_engines 
		WHERE name = $1`, name)

	inferEngine := &model.InferEngine{}
	if err := row.Scan(
		&inferEngine.Name,
		&inferEngine.Framework,
		&inferEngine.Image,
		&inferEngine.Status,
		&inferEngine.CreatedAt,
		&inferEngine.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return inferEngine, nil
}
