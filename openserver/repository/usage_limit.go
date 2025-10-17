package repository

import (
	"context"
	"errors"
	"openserver/model"

	"github.com/jackc/pgx/v5"
)

type UsageLimitRepo struct{}

func UsageLimit() *UsageLimitRepo {
	return &UsageLimitRepo{}
}

func (r *UsageLimitRepo) GetByID(ctx context.Context, workspaceID, modelName string) (*model.UsageLimit, error) {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	usageLimit := &model.UsageLimit{}

	row := conn.QueryRow(ctx, `SELECT workspace_id, model_name, request_limit, token_limit, created_at, updated_at 
		FROM usage_limits 
		WHERE workspace_id=$1 AND model_name=$2`, workspaceID, modelName)
	if err := row.Scan(&usageLimit.WorkspaceID, &usageLimit.ModelName, &usageLimit.RequestLimit, &usageLimit.TokenLimit, &usageLimit.CreatedAt, &usageLimit.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return usageLimit, nil
}

func (r *UsageLimitRepo) ListByWorkspaceID(ctx context.Context, workspaceID string) ([]*model.UsageLimit, error) {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT workspace_id, model_name, request_limit, token_limit, created_at, updated_at
		FROM usage_limits
		WHERE workspace_id = $1
		ORDER BY created_at ASC
	`, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usageLimits := []*model.UsageLimit{}
	for rows.Next() {
		usageLimit := &model.UsageLimit{}
		err := rows.Scan(
			&usageLimit.WorkspaceID,
			&usageLimit.ModelName,
			&usageLimit.RequestLimit,
			&usageLimit.TokenLimit,
			&usageLimit.CreatedAt,
			&usageLimit.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		usageLimits = append(usageLimits, usageLimit)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return usageLimits, nil
}

func (r *UsageLimitRepo) Create(ctx context.Context, usageLimit *model.UsageLimit) error {

	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `INSERT INTO usage_limits (workspace_id, model_name, request_limit, token_limit) VALUES ($1, $2, $3, $4)`,
		usageLimit.WorkspaceID,
		usageLimit.ModelName,
		usageLimit.RequestLimit,
		usageLimit.TokenLimit)

	return err
}

func (r *UsageLimitRepo) Update(ctx context.Context, usageLimit *model.UsageLimit) error {

	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	sql := "UPDATE usage_limits SET request_limit=$1, token_limit=$2, updated_at=NOW() WHERE workspace_id=$3 AND model_name=$4"
	_, err = conn.Exec(ctx, sql, usageLimit.RequestLimit, usageLimit.TokenLimit, usageLimit.WorkspaceID, usageLimit.ModelName)

	return err
}

func (r *UsageLimitRepo) Delete(ctx context.Context, workspaceID, modelName string) error {

	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "DELETE FROM usage_limits WHERE workspace_id=$1 AND model_name=$2", workspaceID, modelName)

	return err
}
