package entity

import (
	"time"
)

type Chat struct {
	ID            string    `bson:"id"`
	ChatName      string    `bson:"chat_name"`
	Users         []string  `bson:"users"`
	IsGroup       bool      `bson:"is_group"`
	GroupAdmin    string    `bson:"group_admin"`
	LatestMessage Message   `bson:"latest_message"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at"`
}
