package model

import (
	"time"
)

type ApiService struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	TopoID    uint64    `json:"topoID"`
	UpdatedAt time.Time `json:"updateAt"`
	CreatedAt time.Time `json:"createAt"`
}
