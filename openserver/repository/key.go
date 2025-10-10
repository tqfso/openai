package repository

import (
	"context"
	"fmt"
	"openserver/model"
	"strings"
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
		return nil, err
	}
	return apiKey, nil
}

func (r *ApiKeyRepo) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]*model.ApiKey, int, error) {
	conn, err := GetPool().Acquire(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer conn.Release()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	where := []string{fmt.Sprintf("user_id=%s", userID)}

	// 查询总数
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM api_keys WHERE %s", strings.Join(where, " AND "))
	var total int
	if err := conn.QueryRow(ctx, countSQL).Scan(&total); err != nil {
		return nil, 0, err
	}

	// 查询分页数据（LIMIT 和 OFFSET 用 fmt.Sprintf 拼接）
	sql := fmt.Sprintf(`
        SELECT id, workspace_id, description, expires_at, created_at, updated_at
        FROM api_keys
        WHERE %s
        ORDER BY created_at DESC
        LIMIT %d OFFSET %d
    `, strings.Join(where, " AND "), pageSize, offset)

	rows, err := conn.Query(ctx, sql)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []*model.ApiKey
	for rows.Next() {
		apiKey := &model.ApiKey{}
		if err := rows.Scan(
			&apiKey.ID,
			&apiKey.WorkspaceID,
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
