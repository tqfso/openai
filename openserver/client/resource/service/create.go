package service

import (
	"common/types"
	"context"
	"openserver/client/resource"
)

type CreateRequest struct {
	User          string           `json:"user" validate:"required,min=34,max=40"`
	Image         string           `json:"image" validate:"required"`
	GpuModel      string           `json:"gpuModel,omitempty"`
	GpuCount      *types.Quantity  `json:"gpuCount,omitempty"`
	CudaMin       string           `json:"cudaMin,omitempty" validate:"field_validator"`
	CudaMax       string           `json:"cudaMax,omitempty" validate:"field_validator"`
	CpuCores      *types.Quantity  `json:"cpuCores,omitempty"`
	MemSizeLimit  *types.Quantity  `json:"memSizeLimit,omitempty"`
	ShmSizeLimit  *types.Quantity  `json:"shmSizeLimit,omitempty"`
	SysSizeLimit  *types.Quantity  `json:"sysSizeLimit,omitempty"`
	DataSizeLimit *types.Quantity  `json:"dataSizeLimit,omitempty"`
	Mounts        []PathMount      `json:"mounts" validate:"dive"`
	BandWidth     *types.BandWidth `json:"bandWidth,omitempty"`
	ReplicaCount  int32            `json:"replicaCount,omitempty" validate:"gte=0,lte=100"`
	Env           []EnvVar         `json:"env,omitempty"`
	Command       []string         `json:"command,omitempty"`
	Args          []string         `json:"args,omitempty"`
	Port          int32            `json:"port,omitempty"`
	AccessMode    string           `json:"accessMode,omitempty" validate:"required,oneof=Ingress Direct Service Balance"`
	Privileged    bool             `json:"privileged,omitempty"`
	ApplyVip      bool             `json:"applyVip,omitempty"`
	EipInfo       *EipInfo         `json:"eipInfo,omitempty"`
	VpcId         types.VpcID      `json:"vpcId,omitempty"`
}

type ServiceID struct {
	Key string `json:"key,omitempty"`
}

func Create(ctx context.Context, req *CreateRequest) (string, error) {
	var resp ServiceID
	if err := resource.Post(ctx, "v1/service/precreate", req, &resp); err != nil {
		return "", err
	}

	serviceID := resp.Key
	if err := resource.Post(ctx, "v1/service/create", resp, nil); err != nil {
		return "", err
	}

	return serviceID, nil
}

func Release(ctx context.Context, id string) error {
	var reqest ServiceID
	if err := resource.Post(ctx, "v1/service/release", reqest, nil); err != nil {
		return err
	}
	return nil
}
