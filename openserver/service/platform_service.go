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

type ServiceTargetResponse struct {
	ID      string
	Err     error
	Targets []*model.ModelServiceTarget
}

// 获取API网关负责的模型服务列表
func (s *PlatformServiceDefault) ListByGateway(ctx context.Context, apiServiceID string) ([]*model.ModelServiceInfo, error) {
	services, err := repository.PlatormService().ListByGateway(ctx, apiServiceID)
	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, nil
	}

	var infoList []*model.ModelServiceInfo
	for _, service := range services {
		infoList = append(infoList, &model.ModelServiceInfo{
			ID:        service.ID,
			ModelName: service.ModelName,
			Power:     service.Power,
			Load:      service.Load,
		})
	}

	ackCount := 0
	channel := make(chan ServiceTargetResponse, len(services))
	for _, service := range services {
		go s.getServiceTargets(ctx, service.ID, channel)
	}

	for {
		response, ok := <-channel
		if !ok {
			return nil, fmt.Errorf("channel error")
		}

		if response.Err != nil {
			return nil, response.Err
		}

		for _, info := range infoList {
			if info.ID == response.ID {
				info.Targets = response.Targets
				break
			}
		}

		ackCount++

		if ackCount == len(services) {
			break
		}
	}

	return infoList, nil
}

func (s *PlatformServiceDefault) getServiceTargets(ctx context.Context, id string, channel chan<- ServiceTargetResponse) {
	status, err := service.GetStatus(ctx, id)
	if err != nil {
		channel <- ServiceTargetResponse{ID: id, Err: err}
		return
	}

	var targets []*model.ModelServiceTarget
	if status.EipInfo == nil {
		for _, target := range status.Replicas {
			targets = append(targets, &model.ModelServiceTarget{Port: 8000, IP: target.PirvateIP})
		}
	} else {
		targets = append(targets, &model.ModelServiceTarget{Port: 8000, IP: status.EipInfo.GetIP()})
	}

	channel <- ServiceTargetResponse{ID: id, Targets: targets}
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
