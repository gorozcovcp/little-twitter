package model

import "time"

type Tweet struct {
	UserID  string    `json:"user_id" bson:"user_id"`
	Content string    `json:"content" bson:"content"`
	Created time.Time `json:"created" bson:"created"`
}
