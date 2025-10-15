package model

import "time"

type PlatformModel struct {
	Name             string      `json:"name"`
	Provider         uint64      `json:"provider"`
	Classes          []uint64    `json:"classes,omitempty"`
	Ability          []uint64    `json:"ability,omitempty"`
	MaxContextLength uint64      `json:"maxContextLength"`
	DeployInfo       *DeployInfo `json:"deployInfo,omitempty"`
	Status           string      `json:"status"`
	UpdatedAt        time.Time   `json:"updateAt"`
	CreatedAt        time.Time   `json:"createAt"`
}

type InferInfo struct {
	Env     []string `json:"envs,omitempty"`    // 环境变量如: ENV1=10...
	Command []string `json:"command,omitempty"` // 命令
	Args    []string `json:"args,omitempty"`    // 参数
}

type SuitableGpu struct {
	Name  string `json:"name"`  // GPU型号
	Count int    `json:"count"` // GPU个数
}

// 部署信息
type DeployInfo struct {
	SuitableGpus []SuitableGpu         `json:"suitableGpus,omitempty"` // 合适的加速卡
	InferEngines map[string]*InferInfo `json:"inferEngines,omitempty"` // 推理引擎对应的模型部署参数
}
