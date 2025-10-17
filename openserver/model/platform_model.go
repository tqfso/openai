package model

import (
	"common/types"
	"openserver/client/resource/service"
	"time"
)

type PlatformModel struct {
	Name             string      `json:"name"`
	Provider         uint64      `json:"provider"`
	Classes          []uint64    `json:"classes,omitempty"`
	Abilities        []uint64    `json:"abilities,omitempty"`
	MaxContextLength uint64      `json:"maxContextLength"`
	DeployInfo       *DeployInfo `json:"deployInfo,omitempty"`
	Description      string      `json:"description,omitempty"`
	Status           string      `json:"status"`
	UpdatedAt        time.Time   `json:"updateAt"`
	CreatedAt        time.Time   `json:"createAt"`
}

// 搜索参数
type PlatformModelSearchParam struct {
	ClassesAny   []uint64 `form:"classesAny"`
	AbilitiesAll []uint64 `form:"abilitiesAll"`
	MinContext   *uint64  `form:"minContext"`
	MaxContext   *uint64  `form:"maxContext"`
	PageIndex    int      `form:"pageIndex"`
	PageSize     int      `form:"pageSize"`
}

type InferInfo struct {
	Name         string              `json:"name" binding:"required"`                   // 推理引擎
	CpuCores     *types.Quantity     `json:"cpuCores,omitempty"`                        // CPU核数
	MemSizeLimit *types.Quantity     `json:"memSizeLimit,omitempty"`                    // 内存大小
	ShmSizeLimit *types.Quantity     `json:"shmSizeLimit,omitempty"`                    // 共享内存大小
	Mounts       []service.PathMount `json:"mounts,omitempty"`                          // 模型挂载目录
	Env          []service.EnvVar    `json:"env,omitempty"`                             // 环境变量
	Command      []string            `json:"command,omitempty"`                         // 命令
	Args         []string            `json:"args,omitempty"`                            // 参数
	SuitableGpus []*SuitableGpu      `json:"suitableGpus,omitempty" binding:"required"` // 合适的AI加速卡
}

type SuitableGpu struct {
	Name    string           `json:"name" binding:"required"`
	Count   int              `json:"count" binding:"required"`
	Env     []service.EnvVar `json:"env,omitempty"`
	Command []string         `json:"command,omitempty"`
	Args    []string         `json:"args,omitempty"`
}

// 部署信息
type DeployInfo struct {
	InferInfos []*InferInfo `json:"inferInfos,omitempty" binding:"required"` // 合适的推理引擎
}

func (info *InferInfo) GetPlatformInferGpu() *SuitableGpu {
	if len(info.SuitableGpus) == 0 {
		return nil
	}
	return info.SuitableGpus[0]
}

func (info *DeployInfo) GetPlatformInferInfo() *InferInfo {
	if len(info.InferInfos) == 0 {
		return nil
	}
	return info.InferInfos[0]
}
