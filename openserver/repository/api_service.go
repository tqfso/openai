package repository

import (
	"context"
	"errors"
	"net"
	"openserver/model"

	"github.com/jackc/pgx/v5"
)

type ApiServiceRepo struct{}

func ApiService() *ApiServiceRepo {
	return &ApiServiceRepo{}
}

func (r *ApiServiceRepo) GetByTopoID(ctx context.Context, topoID string) (*model.ApiService, error) {
	conn, err := GetPool().Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, `
		SELECT id, topo_id, public_ip, created_at, updated_at 
		FROM api_services 
		WHERE topo_id = $1`, topoID)

	apiService := &model.ApiService{}
	if err := row.Scan(
		&apiService.ID,
		&apiService.TopoID,
		&apiService.PublicIP,
		&apiService.CreatedAt,
		&apiService.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return apiService, nil
}

func (r *ApiServiceRepo) Create(ctx context.Context, topoID uint64, id string, ip net.IP) error {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `INSERT INTO api_services (id, topo_id, public_ip,) VALUES ($1, $2, $3)`, id, topoID, ip)

	return err
}

func (r *ApiServiceRepo) Delete(ctx context.Context, id string) error {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `DELETE FROM api_services WHERE id=$1`, id)
	return err
}
