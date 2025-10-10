package repository

import (
	"context"
	"fmt"
	"openserver/model"
	"strings"
)

type WorkspaceRepo struct{}

func Workspace() *WorkspaceRepo {
	return &WorkspaceRepo{}
}

func (r *WorkspaceRepo) GetByID(ctx context.Context, id uint64) (*model.Workspace, error) {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, `SELECT id, user_id, name, status, created_at, updated_at FROM workspaces WHERE id=$1`, id)
	workspace := &model.Workspace{}
	if err := row.Scan(&workspace.ID, &workspace.UserID, &workspace.Name, &workspace.Status, &workspace.CreatedAt, &workspace.UpdatedAt); err != nil {
		return nil, err
	}

	return workspace, nil
}

func (r *WorkspaceRepo) GetCountByUser(ctx context.Context, userID string) (int, error) {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	var total int
	err = conn.QueryRow(ctx, `SELECT COUNT(*) FROM workspaces WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *WorkspaceRepo) ListByUser(ctx context.Context, userID string) ([]*model.Workspace, error) {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT id, user_id, name, status, created_at, updated_at
		FROM workspaces
		WHERE user_id = $1
		ORDER BY created_at ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workspaces := []*model.Workspace{}
	for rows.Next() {
		workspace := &model.Workspace{}
		err := rows.Scan(
			&workspace.ID,
			&workspace.UserID,
			&workspace.Name,
			&workspace.Status,
			&workspace.CreatedAt,
			&workspace.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, workspace)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return workspaces, nil
}

func (r *WorkspaceRepo) Create(ctx context.Context, workspace *model.Workspace) (uint64, error) {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	fieldMap := map[string]any{
		"id":      workspace.ID,
		"user_id": workspace.UserID,
		"name":    workspace.Name,
		"status":  workspace.Status,
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

	id := uint64(0)
	sql := fmt.Sprintf("INSERT INTO workspaces (%s) VALUES (%s) RETURNING id",
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	err = conn.QueryRow(ctx, sql, args...).Scan(&id)
	return id, err
}

func (r *WorkspaceRepo) Delete(ctx context.Context, id uint64, userID string) error {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `DELETE FROM workspaces WHERE id=$1 AND user_id=$2`, id, userID)
	return err
}
