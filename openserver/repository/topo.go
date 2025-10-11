package repository

import (
	"context"
	"errors"
	"openserver/model"

	"github.com/jackc/pgx/v5"
)

type TopoRepo struct{}

func Topo() *TopoRepo {
	return &TopoRepo{}
}

func (r *TopoRepo) GetByID(ctx context.Context, id uint64) (*model.Topo, error) {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, `SELECT id, status, created_at, updated_at FROM topo_domains WHERE id=$1`, id)
	topo := &model.Topo{}
	if err := row.Scan(&topo.ID, &topo.Status, &topo.CreatedAt, &topo.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return topo, nil
}

func (r *TopoRepo) Add(ctx context.Context, topo *model.Topo) error {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `INSERT INTO topo_domains (id, vpc_id) VALUES ($1, $2)`, topo.ID, topo.VpcID)

	return err
}
