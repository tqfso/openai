package service

import (
	"common/types"
	"context"
	"net"
	"openserver/client/resource"
)

type StatusRequest struct {
	Key string `form:"key"`
}

type StatusResponse struct {
	Status         string                 `json:"status,omitempty"`
	VpcId          types.VpcID            `json:"vpcId,omitempty"`
	ExportPort     uint16                 `json:"exportPort,omitempty"`
	AccessDomain   string                 `json:"accessDomain,omitempty"`
	AccessPort     uint16                 `json:"accessPort,omitempty"`
	AccessProtocol string                 `json:"accessProtocol,omitempty"`
	Replicas       map[int]*ReplicaStatus `json:"replicas,omitempty"`
	Service        *Service               `json:"service,omitempty"`
	Balance        *Balance               `json:"balance,omitempty"`
	EipInfo        *EipInfo               `json:"eipInfo,omitempty"`
	Env            []EnvVar               `json:"env,omitempty"`
	Mounts         []PathMount            `json:"mounts,omitempty"`
	Command        []string               `json:"command,omitempty"`
	Args           []string               `json:"args,omitempty"`
}

type ReplicaStatus struct {
	Node      string `json:"node,omitempty"`
	Status    string `json:"status,omitempty"`
	PirvateIP string `json:"privateIP,omitempty"`
	HostIP    string `json:"hostIP,omitempty"`
}

type Service struct {
	BalanceId string         `json:"balanceId,omitempty"` // 负载均衡器ID
	PortList  []*ServicePort `json:"portList,omitempty"`  // 服务端口信息
}

type ServicePort struct {
	BalancePort   uint16 `json:"balancePort,omitempty"`   // 负载均衡端口
	ContainerPort uint16 `json:"containerPort,omitempty"` // 容器服务端口
	TcpProtocol   bool   `json:"tcpProtocol,omitempty"`   // 是否支持TCP
	UdpProtocol   bool   `json:"udpProtocol,omitempty"`   // 是否支持UDP
	PublicPort    uint16 `json:"publicPort,omitempty"`    // 外部访问端口
}

type Balance struct {
	Services map[string]*BalanceService `json:"serives,omitempty"`
}

type BalanceService struct {
	PortList []*ServicePort `json:"portList,omitempty"`
}

type EipInfo struct {
	Gateway      string           `json:"gateway,omitempty" binding:"required"`                        // 网关地址
	IP           string           `json:"ip,omitempty" binding:"required"`                             // 弹性公网地址(CIDR格式)
	FlowStatsURL string           `json:"flowStatsURL,omitempty" binding:"omitempty,url"`              // 流量上报地址
	Status       string           `json:"status,omitempty" binding:"omitempty,oneof=Enabled Disabled"` // 状态
	BandWidth    *types.BandWidth `json:"bandWidth,omitempty"`                                         // 带宽峰值限制
}

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

type PathMount struct {
	Name          string `json:"name,omitempty"`
	HostPath      string `json:"hostPath"`
	ContainerPath string `json:"containerPath"`
	ReadOnly      bool   `json:"readOnly"`
}

func (eip EipInfo) GetIP() string {
	ip, _, err := net.ParseCIDR(eip.IP)
	if err != nil {
		return ""
	}

	return ip.String()
}

func GetStatus(ctx context.Context, id string) (*StatusResponse, error) {

	param := ServiceID{Key: id}
	var resp StatusResponse
	if err := resource.Get(ctx, "/v1/service/status", param, &resp); err != nil {
		return nil, err
	}
	return &resp, nil

}
