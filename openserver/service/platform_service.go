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
)

type PlatformServiceDefault struct {
	EipInfo *service.EipInfo // 调试
}

func PlatformService() *PlatformServiceDefault {
	return &PlatformServiceDefault{}
}

// 获取API网关负责的模型服务列表
func (s *PlatformServiceDefault) ListByGateway(ctx context.Context, apiServiceID string) ([]*model.PlatformService, error) {
	return repository.PlatormService().ListByGateway(ctx, apiServiceID)
}

func (s *PlatformServiceDefault) Create(ctx context.Context, topoID uint64, serviceName, modelName string) (string, error) {

	// 获取拓扑域AIP网关

	apiService, err := Api().FindByTopoID(ctx, topoID)
	if err != nil {
		return "", err
	}

	if apiService == nil {
		return "", &common.Error{Code: common.ApiServiceNotFound, Msg: "api gateway not found"}
	}

	// 创建平台模型服务

	serviceID, err := s.createService(ctx, topoID, modelName)
	if err != nil {
		return "", err
	}

	// 保存平台模型服务

	platformService := &model.PlatformService{
		ID:           serviceID,
		TopoID:       topoID,
		Name:         serviceName,
		ModelName:    modelName,
		ApiServiceID: apiService.ID,
	}

	if err := repository.PlatormService().Create(ctx, platformService); err != nil {
		return "", err
	}

	return serviceID, nil
}

func (s *PlatformServiceDefault) createService(ctx context.Context, topoID uint64, modelName string) (string, error) {

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
		EipInfo:      s.EipInfo,
	}

	request.Mounts = append(request.Mounts, inferInfo.Mounts...)
	if len(request.Mounts) == 0 {
		request.Mounts = []service.PathMount{
			{
				Name:          "model",
				HostPath:      fmt.Sprintf("/mnt/cephfs/openai/models/%s", modelName),
				ContainerPath: fmt.Sprintf("/models/%s", modelName),
				ReadOnly:      true,
			},
		}
	}

	request.Env = append(request.Env, inferInfo.Env...)
	request.Env = append(request.Env, inferGpu.Env...)
	request.Command = append(request.Command, inferInfo.Command...)
	request.Command = append(request.Command, inferGpu.Command...)
	request.Args = append(request.Args, inferInfo.Args...)
	request.Args = append(request.Args, inferGpu.Args...)

	serviceID, err := service.Create(ctx, &request)
	if err != nil {
		return "", err
	}

	return serviceID, nil

}

func (s *PlatformServiceDefault) Release(ctx context.Context, id string) error {
	if err := service.Release(ctx, id); err != nil {
		return err
	}
	return repository.PlatormService().Delete(ctx, id)
}
