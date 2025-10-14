package resource

import (
	"common/types"
	"context"
	"openserver/config"
)

type CreateVPCResponse struct {
	VpcID types.VpcID `json:"vpcId"`
	IP    string      `json:"ip"`
}

func CreateVPC(ctx context.Context, topoID uint32) (uint64, error) {
	request := map[string]any{
		"topId":  topoID,
		"userId": config.GetZdan().CloudUserId,
	}

	response := CreateVPCResponse{}
	if err := Post(ctx, "v1/vpc/create", request, &response); err != nil {
		return 0, err
	}

	return response.VpcID.Number(), nil
}
