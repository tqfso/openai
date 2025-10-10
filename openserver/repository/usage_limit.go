package repository

import (
	"context"
	"openserver/model"
)

type UsageLimitRepo struct{}

func UsageLimit() *UsageLimitRepo {
	return &UsageLimitRepo{}
}

func (r *UsageLimitRepo) Create(ctx context.Context, usageLimit *model.UsageLimit) error {

	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `INSERT INTO usage_limits (workspace_id, service_id, request_limit, token_limit) VALUES ($1, $2, $3, $4)`,
		usageLimit.WorkspaceID,
		usageLimit.ServiceID,
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

	sql := "UPDATE usage_limits SET request_limit=$1, token_limit=$2, updated_at=NOW() WHERE workspace_id=$3 AND service_id=$4"
	_, err = conn.Exec(ctx, sql, usageLimit.RequestLimit, usageLimit.TokenLimit, usageLimit.WorkspaceID, usageLimit.ServiceID)

	return err
}

func (r *UsageLimitRepo) Delete(ctx context.Context, workspaceID uint64, serviceID string) error {

	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "DELETE FROM usage_limits WHERE workspace_id=$1 AND service_id=$2", workspaceID, serviceID)

	return err
}
