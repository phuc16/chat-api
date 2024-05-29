package entity

import "time"

type Profile struct {
	UserID     string    `bson:"user_id" json:"userID"`
	UserName   string    `bson:"user_name" json:"userName"`
	Gender     bool      `bson:"gender" json:"gender"`
	Birthday   time.Time `bson:"birthday" json:"birthday,omitempty"`
	Avatar     string    `bson:"avatar" json:"avatar"`
	Background string    `bson:"background" json:"background"`
	CreatedAt  time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updatedAt"`
}
