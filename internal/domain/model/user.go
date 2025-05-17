package model

type User struct {
	ID      string   `json:"id"`
	Follows []string `json:"follows"`
}
