package dto

import "time"

type UserMessageDTO struct {
	ID  string          `json:"id"`
	TUM TypeUserMessage `json:"tum"` // type user message
}

type AppendConversationDTO struct {
	UserMessageDTO
	SenderID       string `json:"senderID"`
	SenderName     string `json:"senderName"`
	SenderAvatar   string `json:"senderAvatar"`
	ReceiverID     string `json:"receiverID"`
	ReceiverName   string `json:"receiverName"`
	ReceiverAvatar string `json:"receiverAvatar"`
	Type           string `json:"type"` // type of relationship
}

type FriendRequestAddDTO struct {
	UserMessageDTO
	SenderID       string    `json:"senderID"`
	SenderName     string    `json:"senderName"`
	SenderAvatar   string    `json:"senderAvatar"`
	ReceiverID     string    `json:"receiverID"`
	ReceiverName   string    `json:"receiverName"`
	ReceiverAvatar string    `json:"receiverAvatar"`
	Description    string    `json:"description"`
	SendAt         time.Time `json:"sendAt"`
}

type FriendRequestAcceptDTO struct {
	UserMessageDTO
	SenderID       string `json:"senderID"`
	SenderName     string `json:"senderName"`
	SenderAvatar   string `json:"senderAvatar"`
	ReceiverID     string `json:"receiverID"`
	ReceiverName   string `json:"receiverName"`
	ReceiverAvatar string `json:"receiverAvatar"`
}

type FriendRequestRemoveDTO struct {
	UserMessageDTO
	SenderID   string `json:"senderID"`
	ReceiverID string `json:"receiverID"`
}

type UnfriendDTO struct {
	UserMessageDTO
	SenderID   string `json:"senderID"`
	ReceiverID string `json:"receiverID"`
}
