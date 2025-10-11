package service

import (
	"common"
	"context"
	"errors"
	"openserver/model"
	"openserver/repository"

	"github.com/jackc/pgx/v5"
)

type UserService struct{}

func User() *UserService {
	return &UserService{}
}

// 查找用户信息
func (s *UserService) FindByID(ctx context.Context, id string) (*model.User, error) {
	user, err := repository.User().GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = &common.Error{Code: common.UserNotFound, Msg: "user not found"}
		}
		return nil, err
	}

	return user, nil
}

// 用户开通模型开放平台
func (s *UserService) Create(ctx context.Context, id, nickName string, requestLimit, tokenLimit int64) error {
	user := &model.User{
		ID:           id,
		NickName:     nickName,
		RequestLimit: requestLimit,
		TokenLimit:   tokenLimit,
	}

	return repository.User().CreateWithDefaultWorkspace(ctx, user)
}
