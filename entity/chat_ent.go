package entity

import "time"

type Chat struct {
	ID             string         `bson:"id" json:"id"`
	Deliveries     []Delivery     `bson:"deliveries" json:"deliveries"`
	Reads          []Delivery     `bson:"reads" json:"reads"`
	ChatActivities []ChatActivity `bson:"chat_activities" json:"chatActivities"`
	CreatedAt      time.Time      `bson:"created_at" json:"createdAt"`
	UpdatedAt      time.Time      `bson:"updated_at" json:"updatedAt"`
}
