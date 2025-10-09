package service

import (
	"context"
	"openserver/model"
	"openserver/repository"
)

type WorkspaceService struct{}

func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{}
}

// 获取用户工作空间列表
func (s *WorkspaceService) ListAll(ctx context.Context, userID string) ([]*model.Workspace, error) {
	return repository.NewWorkspaceRepo().ListAll(ctx, userID)
}

// 创建工作空间
func (s *WorkspaceService) Create(ctx context.Context, userID, name string) (uint64, error) {
	workspace := model.Workspace{
		UserID: userID,
		Name:   name,
	}

	return repository.NewWorkspaceRepo().Create(ctx, &workspace)
}

// 删除工作空间
func (s *WorkspaceService) Delete(ctx context.Context, id, userID string) error {
	return repository.NewWorkspaceRepo().Delete(ctx, id, userID)
}

// 授权模型服务
func (s *WorkspaceService) CreateUsageLimit(ctx context.Context, workspaceId uint64, serviceId string) error {

	usageLimit := model.UsageLimit{
		WorkspaceID: workspaceId,
		ServiceID:   serviceId,
	}
	return repository.NewUsageLimitRepo().Create(ctx, &usageLimit)
}

// 设置调用限制
func (s *WorkspaceService) UpdateUsageLimit(ctx context.Context, workspaceId uint64, serviceId string, requestLimit, tokenLimit int64) error {
	usageLimit := model.UsageLimit{
		WorkspaceID:  workspaceId,
		ServiceID:    serviceId,
		RequestLimit: requestLimit,
		TokenLimit:   tokenLimit,
	}
	return repository.NewUsageLimitRepo().Update(ctx, &usageLimit)
}
