package entity

import "time"

type PersonInfo struct {
	UserID     string    `bson:"user_id" json:"userID"`
	UserName   string    `bson:"user_name" json:"userName"`
	UserAvatar string    `bson:"user_avatar" json:"userAvatar"`
	CreatedAt  time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updatedAt"`
}
