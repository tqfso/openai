package model

import "time"

type PlatformService struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	TopoID       uint64    `json:"topoID"`
	ModelName    string    `json:"modelName"`
	ApiServiceID string    `json:"apiServiceID"`
	Power        uint64    `json:"power"`
	Load         uint64    `json:"load"`
	UpdatedAt    time.Time `json:"updateAt"`
	CreatedAt    time.Time `json:"createAt"`
}

type ModelServiceInfo struct {
	ID        string                `json:"id"`
	ModelName string                `json:"modelName"`
	Power     uint64                `json:"power"`
	Load      uint64                `json:"load"`
	Targets   []*ModelServiceTarget `json:"targets"`
}

type ModelServiceTarget struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}
