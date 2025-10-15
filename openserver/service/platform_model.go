package service

import (
	"context"
	"openserver/model"
	"openserver/repository"
)

type PlatformModelService struct{}

func PlatformModel() *PlatformModelService {
	return &PlatformModelService{}
}

// 查询指定模型
func (s *PlatformModelService) FindByModelName(ctx context.Context, modelName string) (*model.PlatformModel, error) {
	return repository.PlatformModel().GetByModelName(ctx, modelName)
}

// 查询模型列表
func (r *PlatformModelService) List(ctx context.Context, searchParams *model.PlatformModelSearchParam) ([]*model.PlatformModel, int, error) {
	return repository.PlatformModel().List(ctx, searchParams)
}

// 预置模型
func (r *PlatformModelService) Create(ctx context.Context, pm *model.PlatformModel) error {
	return repository.PlatformModel().Create(ctx, pm)
}

// 删除模型
func (r *PlatformModelService) Delete(ctx context.Context, name string) error {
	return repository.PlatformModel().Delete(ctx, name)
}
