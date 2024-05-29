package entity

type AuthChatMessage struct {
	Token  string `bson:"token" json:"token"`
	Sender string `bson:"sender" json:"sender"`
}
