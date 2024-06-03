package entity

import "time"

type Content struct {
	Key       string    `bson:"key" json:"key"`
	Value     string    `bson:"value" json:"value"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}
