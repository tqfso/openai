package resource

import (
	"context"
	"errors"
	"openserver/config"
)

func CreateVPC(ctx context.Context, topoID uint32) (uint64, error) {
	request := map[string]any{
		"topId":  topoID,
		"userId": config.GetZdan().CloudUserId,
	}

	response := map[string]any{}
	if err := Post(ctx, "v1/vpc/create", request, response); err != nil {
		return 0, err
	}

	vpcId, exists := response["vpcId"]
	if !exists {
		return 0, errors.New("create vpc response data error")
	}

	return vpcId.(uint64), nil
}
