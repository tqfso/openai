package repository

import (
	"context"
	"openserver/model"
)

type TopoRepo struct{}

func NewTopoRepo() *TopoRepo {
	return &TopoRepo{}
}

// 获取某个拓扑域的信息
func (r *TopoRepo) GetByID(ctx context.Context, id uint64) (*model.Topo, error) {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, `SELECT id, status, created_at, updated_at FROM topologies WHERE id=$1`, id)
	topo := &model.Topo{}
	if err := row.Scan(&topo.ID, &topo.Status, &topo.CreatedAt, &topo.UpdatedAt); err != nil {
		return nil, err
	}
	return topo, nil
}

// 新增拓扑域
func (r *TopoRepo) Add(ctx context.Context, topo *model.Topo) error {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "INSERT INTO topo_domains (id, vpc_id) VALUES ($1, $2)", topo.ID, topo.VpcID)

	return err
}
