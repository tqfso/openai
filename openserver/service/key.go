package service

import (
	"common/secure"
	"context"
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

	apiKey, err := repository.ApiKey().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	plainText, err := secure.Decrypt(apiKey.ID)
	if err != nil {
		return nil, err
	}

	apiKey.ID = plainText

	return apiKey, nil
}

// 查询用户密钥列表
func (s *ApiKeyService) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]*model.ApiKey, int, error) {
	apiKeys, totalCount, err := repository.ApiKey().ListByUser(ctx, userID, 0, page, pageSize)
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
func (s *ApiKeyService) Create(ctx context.Context, workspaceID uint64, description string, expiredAt *time.Time) (string, error) {

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
		WorkspaceID: workspaceID,
		Description: description,
		ExpiresAt:   expiredAt,
	}

	if err := repository.ApiKey().Create(ctx, apiKey); err != nil {
		return "", err
	}

	return plainText, nil
}
