package repository

import (
	"context"
	"openserver/model"
)

type PlatformServiceRepo struct{}

func PlatormService() *PlatformServiceRepo {
	return &PlatformServiceRepo{}
}

func (r *PlatformServiceRepo) Create(ctx context.Context, modelService *model.PlatformService) error {

	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `INSERT INTO platform_services (id, name, topo_id, model_name, api_service_id, power, load) 
							VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		modelService.ID,
		modelService.Name,
		modelService.TopoID,
		modelService.ModelName,
		modelService.ApiServiceID,
		modelService.Power,
		modelService.Load,
	)

	return err
}

func (r *PlatformServiceRepo) Delete(ctx context.Context, id string) error {

	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `DELETE FROM platform_services WHERE id=$1`, id)

	return err
}
