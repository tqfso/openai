package service

import (
	"common/types"
	"context"
	"openserver/client/resource/service"
	"openserver/config"
	"openserver/model"
	"openserver/repository"
)

type GatewayService struct{}

func ApiService() *GatewayService {
	return &GatewayService{}
}

func (s *GatewayService) FindByTopoID(ctx context.Context, topoID uint64) (*model.ApiService, error) {
	return repository.ApiService().GetByTopoID(ctx, topoID)
}

func (s *GatewayService) Create(ctx context.Context, topoID uint64, name string, eipInfo *service.EipInfo) (string, error) {

	// 获取私有网络

	vpcID, err := Topo().FetchVpcID(ctx, topoID)
	if err != nil {
		return "", err
	}

	// 创建网关服务

	request := service.CreateRequest{
		User:       config.GetZdan().CloudUserId,
		VpcId:      types.NewVpcId(vpcID),
		AccessMode: "Service",
		Image:      "openai/apiserver:v1.0.0",
		EipInfo:    eipInfo,
	}

	request.Mounts = []service.PathMount{
		{
			Name:          "models",
			HostPath:      "/mnt/cephfs/openai/models",
			ContainerPath: "/models",
		},
	}

	serviceID, err := service.Create(ctx, &request)
	if err != nil {
		return "", err
	}

	// 保存数据库

	apiService := model.ApiService{
		ID:     serviceID,
		TopoID: topoID,
		Name:   name,
	}

	return serviceID, repository.ApiService().Create(ctx, &apiService)
}

func (s *GatewayService) Delete(ctx context.Context, id string) error {
	if err := service.Release(ctx, id); err != nil {
		return err
	}

	return repository.ApiService().Delete(ctx, id)
}
