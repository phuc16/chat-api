package entity

import "time"

type ChatActivity struct {
	ID         string    `bson:"id" json:"id"`
	UserID     string    `bson:"user_id" json:"userID"`
	MessageID  string    `bson:"message_id" json:"messageID"`
	UserAvatar string    `bson:"user_avatar" json:"userAvatar"`
	Timestamp  time.Time `bson:"timestamp" json:"timestamp"`
	ParentID   string    `bson:"parent_id" json:"parentID"`
	Contents   []Content `bson:"contents" json:"contents"`
	Hidden     []string  `bson:"hidden" json:"hidden"`
	Recall     bool      `bson:"recall" json:"recall"`
	CreatedAt  time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updatedAt"`
}
