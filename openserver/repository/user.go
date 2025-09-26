package repository

import (
	"common"
	"context"
	"fmt"
	"openserver/model"
	"strings"
)

// UserRepo 提供用户表相关操作
type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	// 如果用户已存在，返回错误
	if err := conn.QueryRow(ctx, `SELECT id FROM users WHERE id=$1`, user.ID).Scan(&user.ID); err == nil {
		return &common.Error{Code: common.UserExistError, Msg: "User already exists"}
	}

	// 动态拼接字段和参数,只插入非零值字段
	fieldMap := map[string]any{
		"id":            user.ID,
		"nick_name":     user.NickName,
		"request_limit": user.RequestLimit,
		"token_limit":   user.TokenLimit,
	}
	columns := []string{}
	placeholders := []string{}
	args := []any{}
	idx := 1
	for k, v := range fieldMap {
		switch val := v.(type) {
		case string:
			if val == "" {
				continue
			}
		case int, int32, int64:
			if fmt.Sprintf("%v", val) == "0" {
				continue
			}
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

func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx,
		`UPDATE users SET nick_name=$1, request_limit=$2, token_limit=$3, status=$4 WHERE id=$5`,
		user.NickName, user.RequestLimit, user.TokenLimit, user.Status, user.ID,
	)
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
