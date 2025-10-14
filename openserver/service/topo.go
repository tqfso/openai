package service

import (
	"context"
	"openserver/client/resource"
	"openserver/model"
	"openserver/repository"
)

type TopoService struct{}

func Topo() *TopoService {
	return &TopoService{}
}

// 获取指定拓扑域的私有网络ID

func (s *TopoService) FetchVpcID(ctx context.Context, topoID uint64) (uint64, error) {

	topoRepo := repository.Topo()
	topo, err := topoRepo.GetByID(ctx, topoID)
	if err != nil {
		return 0, err
	}

	if topo == nil {

		// 没有对应的记录，创建一个VPC并保存数据库

		vpcId, err := resource.CreateVPC(ctx, uint32(topoID))
		if err != nil {
			return 0, err
		}

		topo = &model.Topo{ID: topoID, VpcID: vpcId}
		if err := topoRepo.Add(ctx, topo); err != nil {
			return 0, err
		}

		return vpcId, nil
	}

	return topo.VpcID, nil
}
