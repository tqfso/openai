package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"openserver/model"
	"strings"

	"github.com/jackc/pgx/v5"
)

type PlatformModelRepo struct{}

func PlatformModel() *PlatformModelRepo {
	return &PlatformModelRepo{}
}

func (r *PlatformModelRepo) GetByModelName(ctx context.Context, modelName string) (*model.PlatformModel, error) {
	conn, err := GetPool().Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	const querySQL = `
		SELECT 
			name, provider, classes, abilities, max_context_length, deploy_info, description, status, created_at, updated_at
		FROM platform_models
		WHERE name = $1
	`

	var platModel model.PlatformModel
	var deployInfoJSON []byte

	err = conn.QueryRow(ctx, querySQL, modelName).Scan(
		&platModel.Name,
		&platModel.Provider,
		&platModel.Classes,
		&platModel.Abilities,
		&platModel.MaxContextLength,
		&deployInfoJSON,
		&platModel.Description,
		&platModel.Status,
		&platModel.CreatedAt,
		&platModel.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// 解析 deploy_info JSONB
	if deployInfoJSON != nil {
		var deployInfo model.DeployInfo
		if err := json.Unmarshal(deployInfoJSON, &deployInfo); err != nil {
			return nil, fmt.Errorf("failed to unmarshal deploy_info: %w", err)
		}
		platModel.DeployInfo = &deployInfo
	}

	return &platModel, nil
}

func (r *PlatformModelRepo) List(ctx context.Context, searchParams *model.PlatformModelSearchParam) ([]*model.PlatformModel, int, error) {
	conn, err := GetPool().Acquire(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer conn.Release()

	var where []string
	var args []any
	argIdx := 1

	if len(searchParams.ClassesAny) > 0 {
		where = append(where, fmt.Sprintf("classes && $%d", argIdx))
		args = append(args, searchParams.ClassesAny)
		argIdx++
	}
	if len(searchParams.AbilitiesAll) > 0 {
		where = append(where, fmt.Sprintf("abilities @> $%d", argIdx))
		args = append(args, searchParams.AbilitiesAll)
		argIdx++
	}
	if searchParams.MinContext != nil {
		where = append(where, fmt.Sprintf("max_context_length >= $%d", argIdx))
		args = append(args, *searchParams.MinContext)
		argIdx++
	}
	if searchParams.MaxContext != nil {
		where = append(where, fmt.Sprintf("max_context_length < $%d", argIdx))
		args = append(args, *searchParams.MaxContext)
		argIdx++
	}

	whereSQL := ""
	if len(where) > 0 {
		whereSQL = " WHERE " + strings.Join(where, " AND ")
	}

	// total count
	countSQL := "SELECT COUNT(*) FROM platform_models" + whereSQL
	var total int
	if err := conn.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// select rows with pagination
	selectSQL := fmt.Sprintf(`
		SELECT
			name, provider, classes, abilities, max_context_length, description, status, created_at, updated_at
		FROM platform_models
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereSQL, argIdx, argIdx+1)

	args = append(args, searchParams.PageSize, (searchParams.PageIndex-1)*searchParams.PageSize)

	rows, err := conn.Query(ctx, selectSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []*model.PlatformModel
	for rows.Next() {
		var pm model.PlatformModel

		if err := rows.Scan(
			&pm.Name,
			&pm.Provider,
			&pm.Classes,
			&pm.Abilities,
			&pm.MaxContextLength,
			&pm.Description,
			&pm.Status,
			&pm.CreatedAt,
			&pm.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		results = append(results, &pm)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (r *PlatformModelRepo) Create(ctx context.Context, pm *model.PlatformModel) error {
	conn, err := GetPool().Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	var deployInfoJSON []byte
	if pm.DeployInfo != nil {
		deployInfoJSON, err = json.Marshal(pm.DeployInfo)
		if err != nil {
			return err
		}
	}

	fieldMap := map[string]any{
		"name":               pm.Name,
		"provider":           pm.Provider,
		"classes":            pm.Classes,
		"abilities":          pm.Abilities,
		"max_context_length": pm.MaxContextLength,
		"deploy_info":        deployInfoJSON,
		"description":        pm.Description,
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

	sql := fmt.Sprintf("INSERT INTO platform_models (%s) VALUES (%s)",
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))
	_, err = conn.Exec(ctx, sql, args...)
	return err
}

func (r *PlatformModelRepo) Delete(ctx context.Context, modelName string) error {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `DELETE FROM platform_models WHERE name=$1`, modelName)
	return err
}
