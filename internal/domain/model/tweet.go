package model

import "time"

type Tweet struct {
	UserID  string    `json:"user_id"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
}
