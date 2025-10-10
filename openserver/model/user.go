package model

import "time"

type User struct {
	ID           string
	NickName     string
	RequestLimit int64
	TokenLimit   int64
	Status       string
	UpdatedAt    time.Time
	CreatedAt    time.Time
}
