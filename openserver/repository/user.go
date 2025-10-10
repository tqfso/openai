package repository

import (
	"common"
	"context"
	"fmt"
	"openserver/model"
	"strings"

	"github.com/jackc/pgx/v5"
)

type UserRepo struct{}

func User() *UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) Exists(ctx context.Context, id string) (bool, error) {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	var exists bool
	err = conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)`, id).Scan(&exists)
	return exists, err
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, `SELECT id, nick_name, request_limit, token_limit, status, created_at, updated_at FROM users WHERE id=$1`, id)
	user := &model.User{}
	if err := row.Scan(&user.ID, &user.NickName, &user.RequestLimit, &user.TokenLimit, &user.Status, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	// 用户是否存在
	exists, err := r.Exists(ctx, user.ID)
	if err != nil {
		return err
	}

	if exists {
		return &common.Error{Code: common.UserExistError, Msg: "user already exists"}
	}

	// 动态拼接字段和参数,只插入非零值字段
	fieldMap := map[string]any{
		"id":            user.ID,
		"nick_name":     user.NickName,
		"request_limit": user.RequestLimit,
		"token_limit":   user.TokenLimit,
		"status":        user.Status,
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

	sql := fmt.Sprintf("INSERT INTO users (%s) VALUES (%s)",
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))
	_, err = conn.Exec(ctx, sql, args...)
	return err
}

func (r *UserRepo) CreateWithDefaultWorkspace(ctx context.Context, user *model.User) error {

	return WithTx(ctx, func(tx pgx.Tx) error {

		var exists bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)`, user.ID).Scan(&exists); err != nil {
			return err
		}

		if exists {
			return &common.Error{Code: common.UserExistError, Msg: "user already exists"}
		}

		// 插入用户记录

		fieldMap := map[string]any{
			"id":            user.ID,
			"nick_name":     user.NickName,
			"request_limit": user.RequestLimit,
			"token_limit":   user.TokenLimit,
			"status":        user.Status,
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

		sql := fmt.Sprintf("INSERT INTO users (%s) VALUES (%s)",
			strings.Join(columns, ", "),
			strings.Join(placeholders, ", "))
		if _, err := tx.Exec(ctx, sql, args...); err != nil {
			return err
		}

		// 插入默认工作空间

		fieldMap = map[string]any{
			"user_id": user.ID,
			"name":    "default",
		}
		columns = []string{}
		placeholders = []string{}
		args = []any{}
		idx = 1
		for k, v := range fieldMap {
			if IsZeroValue(v) {
				continue
			}

			columns = append(columns, k)
			placeholders = append(placeholders, fmt.Sprintf("$%d", idx))
			args = append(args, v)
			idx++
		}

		sql = fmt.Sprintf("INSERT INTO workspaces (%s) VALUES (%s)",
			strings.Join(columns, ", "),
			strings.Join(placeholders, ", "))

		if _, err := tx.Exec(ctx, sql, args...); err != nil {
			return err
		}

		return nil
	})
}

func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	fieldMap := map[string]any{
		"nick_name":     user.NickName,
		"request_limit": user.RequestLimit,
		"token_limit":   user.TokenLimit,
		"status":        user.Status,
	}

	columns := []string{}
	args := []any{}
	idx := 1

	for k, v := range fieldMap {
		if IsZeroValue(v) {
			continue
		}
		columns = append(columns, fmt.Sprintf("%s=$%d", k, idx))
		args = append(args, v)
		idx++
	}

	// 自动更新时间
	columns = append(columns, "updated_at=NOW()")

	sql := fmt.Sprintf("UPDATE users SET %s WHERE id=$%d", strings.Join(columns, ", "), idx)
	args = append(args, user.ID)

	_, err = conn.Exec(ctx, sql, args...)
	return err
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	return err
}
