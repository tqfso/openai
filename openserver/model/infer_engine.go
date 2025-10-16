package model

import "time"

type InferEngine struct {
	Name      string
	Framework string
	Image     string
	Status    string
	UpdatedAt time.Time
	CreatedAt time.Time
}
