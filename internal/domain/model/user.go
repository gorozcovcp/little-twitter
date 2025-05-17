package model

type User struct {
	ID      string   `json:"id" bson:"_id"`
	Follows []string `json:"follows" bson:"follows"`
}
