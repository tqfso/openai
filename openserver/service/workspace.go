package service

import (
	"common"
	"context"
	"fmt"
	"openserver/model"
	"openserver/repository"
)

type WorkspaceService struct{}

func Workspace() *WorkspaceService {
	return &WorkspaceService{}
}

// 查询指定工作空间
func (s *WorkspaceService) FindWorkspace(ctx context.Context, id uint64) (*model.Workspace, error) {
	return repository.Workspace().GetByID(ctx, id)
}

// 查询用户工作空间列表
func (s *WorkspaceService) ListAll(ctx context.Context, userID string) ([]*model.Workspace, error) {
	return repository.Workspace().ListAll(ctx, userID)
}

// 创建工作空间
func (s *WorkspaceService) Create(ctx context.Context, userID, name string) (uint64, error) {

	count, err := repository.Workspace().GetCountByUser(ctx, userID)
	if err != nil {
		return 0, err
	}

	if count >= model.MaxWorkspaceCount {
		return 0, &common.Error{
			Code: common.WorkspaceCountLimit,
			Msg:  fmt.Sprintf("The workspace has reached the maximum number: %d", model.MaxWorkspaceCount),
		}
	}

	workspace := model.Workspace{
		UserID: userID,
		Name:   name,
	}

	return repository.Workspace().Create(ctx, &workspace)
}

// 删除工作空间
func (s *WorkspaceService) Delete(ctx context.Context, id uint64, userID string) error {
	return repository.Workspace().Delete(ctx, id, userID)
}

// 授权模型服务
func (s *WorkspaceService) CreateUsageLimit(ctx context.Context, workspaceId uint64, serviceId string) error {

	usageLimit := model.UsageLimit{
		WorkspaceID: workspaceId,
		ServiceID:   serviceId,
	}
	return repository.UsageLimit().Create(ctx, &usageLimit)
}

// 设置调用限制
func (s *WorkspaceService) UpdateUsageLimit(ctx context.Context, workspaceId uint64, serviceId string, requestLimit, tokenLimit int64) error {
	usageLimit := model.UsageLimit{
		WorkspaceID:  workspaceId,
		ServiceID:    serviceId,
		RequestLimit: requestLimit,
		TokenLimit:   tokenLimit,
	}
	return repository.UsageLimit().Update(ctx, &usageLimit)
}
