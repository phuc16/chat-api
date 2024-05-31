package entity

import "time"

type Delivery struct {
	UserID     string    `bson:"user_id" json:"userID"`
	MessageID  string    `bson:"message_id" json:"messageID"`
	UserAvatar string    `bson:"user_avatar" json:"userAvatar"`
	UserName   string    `bson:"user_name" json:"userName"`
	CreatedAt  time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updatedAt"`
}
