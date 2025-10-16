package service

import (
	"context"
	"openserver/model"
	"openserver/repository"
)

type InferEngineService struct{}

func InferEngine() *InferEngineService {
	return &InferEngineService{}
}

func (s *InferEngineService) FindByName(ctx context.Context, name string) (*model.InferEngine, error) {
	return repository.InferEngine().GetByName(ctx, name)
}
