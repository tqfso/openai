package service

import (
	"context"
	"openserver/model"
	"openserver/repository"
)

type UserService struct{}

func User() *UserService {
	return &UserService{}
}

// 用户开通模型开放平台
func (s *UserService) Create(ctx context.Context, id, nickName string, requestLimit, tokenLimit int64) error {
	user := &model.User{
		ID:           id,
		NickName:     nickName,
		RequestLimit: requestLimit,
		TokenLimit:   tokenLimit,
	}

	return repository.User().Create(ctx, user)
}
