package service

import (
	"common"
	"context"
	"fmt"
	"openserver/model"
	"openserver/repository"

	"github.com/google/uuid"
)

type WorkspaceService struct{}

func Workspace() *WorkspaceService {
	return &WorkspaceService{}
}

// 查询指定工作空间
func (s *WorkspaceService) FindByID(ctx context.Context, id string) (*model.Workspace, error) {
	workspace, err := repository.Workspace().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if workspace == nil {
		return nil, &common.Error{Code: common.WorkspaceNotFound, Msg: "workspace not found"}
	}

	return workspace, nil
}

// 查询用户工作空间列表
func (s *WorkspaceService) ListByUser(ctx context.Context, userID string) ([]*model.Workspace, error) {
	return repository.Workspace().ListByUser(ctx, userID)
}

// 创建工作空间
func (s *WorkspaceService) Create(ctx context.Context, userID, name string) (string, error) {
	count, err := repository.Workspace().GetCountByUser(ctx, userID)
	if err != nil {
		return "", err
	}

	if count >= model.MaxWorkspaceCount {
		return "", &common.Error{Code: common.WorkspaceCountLimit, Msg: fmt.Sprintf("The workspace has reached the maximum number: %d", model.MaxWorkspaceCount)}
	}

	workspace := model.Workspace{
		ID:     uuid.New().String(),
		UserID: userID,
		Name:   name,
	}

	return workspace.ID, repository.Workspace().Create(ctx, &workspace)
}

// 删除工作空间
func (s *WorkspaceService) Delete(ctx context.Context, id string, userID string) error {
	return repository.Workspace().Delete(ctx, id, userID)
}

// 工作空间授权列表
func (s *WorkspaceService) ListUsageLimits(ctx context.Context, workespaceID string) ([]*model.UsageLimit, error) {
	return repository.UsageLimit().ListByWorkspaceID(ctx, workespaceID)
}

// 授权模型服务
func (s *WorkspaceService) CreateUsageLimit(ctx context.Context, workspaceID, serviceID string) error {
	usageLimit := model.UsageLimit{
		WorkspaceID: workspaceID,
		ServiceID:   serviceID,
	}
	return repository.UsageLimit().Create(ctx, &usageLimit)
}

// 设置调用限制
func (s *WorkspaceService) UpdateUsageLimit(ctx context.Context, workspaceID, serviceID string, requestLimit, tokenLimit int64) error {
	usageLimit := model.UsageLimit{
		WorkspaceID:  workspaceID,
		ServiceID:    serviceID,
		RequestLimit: requestLimit,
		TokenLimit:   tokenLimit,
	}
	return repository.UsageLimit().Update(ctx, &usageLimit)
}
