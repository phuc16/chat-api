package entity

import "time"

type FriendRequest struct {
	ID          string    `bson:"id" json:"id"`
	UserID      string    `bson:"user_id" json:"userID"`
	UserName    string    `bson:"user_name" json:"userName"`
	UserAvatar  string    `bson:"user_avatar" json:"userAvatar"`
	Description string    `bson:"description" json:"description"`
	SendAt      time.Time `bson:"send_at" json:"sendAt"`
	IsSender    bool      `bson:"is_sender" json:"isSender"`
	CreatedAt   time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updatedAt"`
}
