package entity

import "time"

type Message struct {
	ID        string    `bson:"id"`
	Sender    string    `bson:"sender"`
	Message   string    `bson:"message"`
	ChatID    string    `bson:"chat_id"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
