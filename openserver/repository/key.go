package repository

import (
	"context"
	"errors"
	"fmt"
	"openserver/model"
	"strings"

	"github.com/jackc/pgx/v5"
)

type ApiKeyRepo struct{}

func ApiKey() *ApiKeyRepo {
	return &ApiKeyRepo{}
}

func (r *ApiKeyRepo) GetByID(ctx context.Context, id string) (*model.ApiKey, error) {
	conn, err := GetPool().Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, `
		SELECT id, user_id, workspace_id, description, expires_at, created_at, updated_at 
		FROM api_keys 
		WHERE id = $1`, id)

	apiKey := &model.ApiKey{}
	if err := row.Scan(
		&apiKey.ID,
		&apiKey.UserID,
		&apiKey.WorkspaceID,
		&apiKey.Description,
		&apiKey.ExpiresAt,
		&apiKey.CreatedAt,
		&apiKey.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return apiKey, nil
}

func (r *ApiKeyRepo) ListByUser(ctx context.Context, userID string, pageIndex, pageSize int) ([]*model.ApiKeyEx, int, error) {
	conn, err := GetPool().Acquire(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer conn.Release()

	offset := (pageIndex - 1) * pageSize

	// 参数化查询，防止SQL注入
	countSQL := `SELECT COUNT(*) FROM api_keys WHERE user_id = $1`
	var total int
	if err := conn.QueryRow(ctx, countSQL, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	querySQL := `
        SELECT a.id, a.user_id, a.workspace_id, w.name AS workspace_name, a.description, a.expires_at, a.created_at, a.updated_at
        FROM api_keys AS a
        JOIN workspaces AS w ON a.workspace_id = w.id
        WHERE a.user_id = $1
        ORDER BY a.created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := conn.Query(ctx, querySQL, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []*model.ApiKeyEx
	for rows.Next() {
		apiKey := &model.ApiKeyEx{}
		if err := rows.Scan(
			&apiKey.ID,
			&apiKey.UserID,
			&apiKey.WorkspaceID,
			&apiKey.WorkspaceName,
			&apiKey.Description,
			&apiKey.ExpiresAt,
			&apiKey.CreatedAt,
			&apiKey.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		results = append(results, apiKey)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (r *ApiKeyRepo) Create(ctx context.Context, apiKey *model.ApiKey) error {
	conn, err := GetPool().Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	fieldMap := map[string]any{
		"id":           apiKey.ID,
		"user_id":      apiKey.UserID,
		"workspace_id": apiKey.WorkspaceID,
		"description":  apiKey.Description,
		"expires_at":   apiKey.ExpiresAt,
	}
	columns := []string{}
	placeholders := []string{}
	args := []any{}
	idx := 1
	for k, v := range fieldMap {
		if IsZeroValue(v) {
			continue
		}
		columns = append(columns, k)
		placeholders = append(placeholders, fmt.Sprintf("$%d", idx))
		args = append(args, v)
		idx++
	}

	sql := fmt.Sprintf("INSERT INTO api_keys (%s) VALUES (%s)",
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))
	_, err = conn.Exec(ctx, sql, args...)
	return err
}

func (r *ApiKeyRepo) Delete(ctx context.Context, id string, userID string) error {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `DELETE FROM api_keys WHERE id=$1 AND user_id=$2`, id, userID)
	return err
}
