package service

import (
	"common"
	"common/secure"
	"context"
	"errors"
	"openserver/model"
	"openserver/repository"
	"time"
)

type ApiKeyService struct{}

func ApiKey() *ApiKeyService {
	return &ApiKeyService{}
}

// 查询指定密钥
func (s *ApiKeyService) FindByID(ctx context.Context, id string) (*model.ApiKey, error) {

	cipherText, err := secure.Encrypt(id)
	if err != nil {
		return nil, err
	}

	apiKey, err := repository.ApiKey().GetByID(ctx, cipherText)
	if err != nil {
		return nil, err
	}

	if apiKey == nil {
		return nil, &common.Error{Code: common.ApiKeyNotFound, Msg: "API KEY not found"}
	}

	apiKey.ID = id

	return apiKey, nil
}

// 查询用户密钥列表
func (s *ApiKeyService) ListByUser(ctx context.Context, userID string, pageIndex, pageSize int) ([]*model.ApiKeyEx, int, error) {
	apiKeys, totalCount, err := repository.ApiKey().ListByUser(ctx, userID, pageIndex, pageSize)
	if err != nil {
		return nil, 0, err
	}

	for _, key := range apiKeys {
		plainText, err := secure.Decrypt(key.ID)
		if err != nil {
			return nil, 0, err
		}

		key.ID = plainText
	}

	return apiKeys, totalCount, nil

}

// 创建密钥
func (s *ApiKeyService) Create(ctx context.Context, userID, workspaceID, description string, expiredAt *time.Time) (string, error) {

	// 判断工作空间是否属于该用户
	workspace, err := Workspace().FindByID(ctx, workspaceID)
	if err != nil {
		return "", err
	}

	if workspace.UserID != userID {
		return "", errors.New("workspace id ownner error")
	}

	// 先随机生成，再加密存储

	plainText, err := secure.GenerateApiKey()
	if err != nil {
		return "", err
	}

	cipherText, err := secure.Encrypt(plainText)
	if err != nil {
		return "", err
	}

	apiKey := &model.ApiKey{
		ID:          cipherText,
		UserID:      userID,
		WorkspaceID: workspaceID,
		Description: description,
		ExpiresAt:   expiredAt,
	}

	if err := repository.ApiKey().Create(ctx, apiKey); err != nil {
		return "", err
	}

	return plainText, nil
}

// 删除密钥
func (s *ApiKeyService) Delete(ctx context.Context, id, userID string) error {

	cipherText, err := secure.Encrypt(id)
	if err != nil {
		return err
	}

	return repository.ApiKey().Delete(ctx, cipherText, userID)
}
