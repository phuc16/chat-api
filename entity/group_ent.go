package entity

import "time"

type Group struct {
	ID         string       `bson:"id" json:"id"`
	ChatName   string       `bson:"chat_name" json:"chatName"`
	ChatAvatar string       `bson:"chat_avatar" json:"chatAvatar"`
	Owner      PersonInfo   `bson:"owner" json:"owner"`
	Admins     []PersonInfo `bson:"admins" json:"admins"`
	Members    []PersonInfo `bson:"members" json:"members"`
	Setting    GroupSetting `bson:"setting" json:"setting"`
	CreatedAt  time.Time    `bson:"created_at" json:"createdAt"`
	UpdatedAt  time.Time    `bson:"updated_at" json:"updatedAt"`
}
