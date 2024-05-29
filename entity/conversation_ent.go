package entity

import "time"

type Conversation struct {
	ID                string         `bson:"id" json:"id"`
	ChatID            string         `bson:"chat_id" json:"chatId"`
	IDUserOrGroup     string         `bson:"id_user_or_group" json:"id_UserOrGroup"`
	ChatName          string         `bson:"chat_name" json:"chatName"`
	ChatAvatar        string         `bson:"chat_avatar" json:"chatAvatar"`
	Type              string         `bson:"type" json:"type"`
	Deliveries        []Delivery     `bson:"deliveries" json:"deliveries"`
	Reads             []Delivery     `bson:"reads" json:"reads"`
	TopChatActivities []ChatActivity `bson:"top_chat_activities" json:"topChatActivities"`
	CreatedAt         time.Time      `bson:"created_at" json:"createdAt"`
	UpdatedAt         time.Time      `bson:"updated_at" json:"updatedAt"`
}
