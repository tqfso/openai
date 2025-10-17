package repository

import (
	"context"
	"openserver/model"
)

type PlatformServiceRepo struct{}

func PlatormService() *PlatformServiceRepo {
	return &PlatformServiceRepo{}
}

func (r *PlatformServiceRepo) ListByGateway(ctx context.Context, apiServiceID string) ([]*model.PlatformService, error) {
	pool := GetPool()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT id, name, topo_id, model_name, api_service_id, power, load, created_at, updated_at
		FROM platform_services
		WHERE api_service_id = $1`, apiServiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	services := []*model.PlatformService{}
	for rows.Next() {
		service := &model.PlatformService{}
		err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.TopoID,
			&service.ModelName,
			&service.ApiServiceID,
			&service.Power,
			&service.Load,
			&service.CreatedAt,
			&service.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
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
