package model

import (
	"net"
	"time"
)

type ApiService struct {
	ID        string    `json:"id"`
	TopoID    uint64    `json:"topoID"`
	PublicIP  net.IP    `json:"publicIP"`
	UpdatedAt time.Time `json:"updateAt"`
	CreatedAt time.Time `json:"createAt"`
}
