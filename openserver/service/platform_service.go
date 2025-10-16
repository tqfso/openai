package service

import (
	"common"
	"common/types"
	"context"
	"fmt"
	"openserver/client/resource/service"
	"openserver/config"
	"openserver/model"
	"openserver/repository"
	"path"
)

type PlatformService struct{}

func Platform() *PlatformService {
	return &PlatformService{}
}

func (s *PlatformService) Create(ctx context.Context, topoID uint64, serviceName, modelName string) error {

	// 获取拓扑域AIP网关

	apiService, err := Api().FindByTopoID(ctx, topoID)
	if err != nil {
		return err
	}

	if apiService == nil {
		return &common.Error{Code: common.ApiServiceNotFound, Msg: "api gateway not found"}
	}

	// 创建平台模型服务

	serviceID, err := s.createService(ctx, topoID, modelName)
	if err != nil {
		return err
	}

	// 保存平台模型服务

	platformService := &model.PlatformService{
		ID:           serviceID,
		TopoID:       topoID,
		ApiServiceID: apiService.ID,
		Name:         serviceName,
		ModelName:    modelName,
	}

	return repository.PlatormService().Create(ctx, platformService)

}

func (s *PlatformService) createService(ctx context.Context, topoID uint64, modelName string) (string, error) {

	// 获取平台预置模型信息

	platormModel, err := PlatformModel().FindByModelName(ctx, modelName)
	if err != nil {
		return "", err
	}

	if platormModel == nil {
		return "", &common.Error{Code: common.PlatModelNotFound, Msg: "platform model not found"}
	}

	deployInfo := platormModel.DeployInfo
	if deployInfo == nil {
		return "", &common.Error{Code: common.PlatModelNotFound, Msg: "platform deploy infomation not found"}
	}

	// 获取部署信息

	inferInfo := deployInfo.GetPlatformInferInfo()
	if inferInfo == nil {
		return "", fmt.Errorf("not found any infer engine")
	}

	inferGpu := inferInfo.GetPlatformInferGpu()
	if inferGpu == nil {
		return "", fmt.Errorf("not found any infer GPU")
	}

	inferEngine, err := InferEngine().FindByName(ctx, inferInfo.Name)
	if err != nil {
		return "", err
	}
	if inferEngine == nil {
		return "", fmt.Errorf("%s not exist", inferInfo.Name)
	}

	// 获取私有网络

	vpcID, err := Topo().FetchVpcID(ctx, topoID)
	if err != nil {
		return "", err
	}

	// 创建平台服务

	request := service.CreateRequest{
		User:         config.GetZdan().CloudUserId,
		VpcId:        types.NewVpcId(vpcID),
		AccessMode:   "Service",
		Image:        inferEngine.Image,
		CpuCores:     inferInfo.CpuCores,
		MemSizeLimit: inferInfo.MemSizeLimit,
		ShmSizeLimit: inferInfo.ShmSizeLimit,
		GpuModel:     inferGpu.Name,
		GpuCount:     types.NewQuantity(int64(inferGpu.Count), types.DecimalExponent),
	}

	if inferInfo.ModelPath == "" {
		inferInfo.ModelPath = "/models"
	}

	request.Mounts = []service.PathMount{
		{
			Name:          "models",
			HostPath:      fmt.Sprintf("/mnt/cephfs/openai/models/%s", modelName),
			ContainerPath: path.Join(inferInfo.ModelPath, modelName),
			ReadOnly:      true,
		},
	}

	serviceID, err := service.Create(ctx, &request)
	if err != nil {
		return "", err
	}

	return serviceID, nil

}
