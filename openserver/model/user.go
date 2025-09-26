package model

type User struct {
	ID           string
	NickName     string
	RequestLimit int64
	TokenLimit   int64
	Status       string
	UpdatedAt    string
	CreatedAt    string
}
