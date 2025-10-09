package service

import (
	"context"
	"openserver/model"
	"openserver/repository"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// 开通用户
func (s *UserService) CreateUser(ctx context.Context, id, nickName string, requestLimit, tokenLimit int64) error {
	user := &model.User{
		ID:           id,
		NickName:     nickName,
		RequestLimit: requestLimit,
		TokenLimit:   tokenLimit,
	}

	return repository.NewUserRepo().Create(ctx, user)
}
