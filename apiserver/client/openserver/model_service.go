package openserver

import (
	"apiserver/config"
	"context"
)

// 查询自己负责的模型服务

type ModelServicesRequest struct {
	ID string `form:"id" binding:"required"`
}

type ModelServicesResponse struct {
	ID        string               `json:"id"`
	ModelName string               `json:"modelName"`
	Power     uint64               `json:"power"`
	Load      uint64               `json:"load"`
	Targets   []ModelServiceTarget `json:"targets"`
}

type ModelServiceTarget struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

func FindModelServices(ctx context.Context) ([]ModelServicesResponse, error) {
	request := ModelServicesRequest{ID: config.GetZdan().ApiServiceId}
	response := []ModelServicesResponse{}
	if err := Get(ctx, "/v1/gateway/model/services", request, &response); err != nil {
		return nil, err
	}

	return response, nil
}
