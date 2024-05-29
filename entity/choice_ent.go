package entity

import "time"

type Choice struct {
	Name      string       `bson:"name" json:"name"`
	Voters    []PersonInfo `bson:"voters" json:"voters"`
	CreatedAt time.Time    `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time    `bson:"updated_at" json:"updatedAt"`
}
