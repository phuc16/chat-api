package entity

import "time"

type Voting struct {
	ID         string     `bson:"id" json:"id"`
	Name       string     `bson:"name" json:"name"`
	Owner      PersonInfo `bson:"owner" json:"owner"`
	DateCreate time.Time  `bson:"date_create" json:"dateCreate"`
	DateLock   time.Time  `bson:"date_lock" json:"dateLock"`
	Choices    []Choice   `bson:"choices" json:"choices"`
	Lock       bool       `bson:"lock" json:"lock"`
	CreatedAt  time.Time  `bson:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `bson:"updated_at" json:"updatedAt"`
}
