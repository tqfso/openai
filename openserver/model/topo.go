package model

import "time"

type Topo struct {
	ID        uint64
	VpcID     uint64
	Status    string
	UpdatedAt time.Time
	CreatedAt time.Time
}
