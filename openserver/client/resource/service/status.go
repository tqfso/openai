package service

import "net"

type StatusRequest struct {
	Key string `form:"key"`
}

type StatusResponse struct {
	Status         string                 `json:"status,omitempty"`
	VpcId          uint64                 `json:"vpcId,omitempty"`
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
	Gateway      net.IP    `json:"gateway,omitempty"`                                            // 网关地址
	IP           net.IPNet `json:"ip,omitempty"`                                                 // 弹性公网地址(CIDR格式)
	FlowStatsURL string    `json:"flowStatsURL,omitempty" validate:"omitempty,url"`              // 流量上报地址
	Status       string    `json:"status,omitempty" validate:"omitempty,oneof=Enabled Disabled"` // 状态
	BandWidth    uint64    `json:"bandWidth,omitempty"`                                          // 带宽峰值限制
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
