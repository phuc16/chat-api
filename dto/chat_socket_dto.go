package dto

import (
	"app/entity"
	"time"
)

type ChatMessageDTO struct {
	ID  string          `json:"id"`
	TCM TypeChatMessage `json:"tcm"` // type chat message
}

type MessageAppendDTO struct {
	ChatMessageDTO
	UserID     string           `json:"userID"`
	UserAvatar string           `json:"userAvatar"`
	UserName   string           `json:"userName"`
	Timestamp  time.Time        `json:"timestamp"`
	ParentID   string           `json:"parentID"`
	Contents   []entity.Content `json:"contents"`
}

type MessageDeliveryDTO struct {
	ChatMessageDTO
	UserID     string `json:"userID"`
	MessageID  string `json:"messageID"`
	UserAvatar string `json:"userAvatar"`
	UserName   string `json:"userName"`
}

type MessageHiddenDTO struct {
	ChatMessageDTO
	UserID    string `json:"userID"`
	MessageID string `json:"messageID"`
}

type AppendVoterDTO struct {
	ChatMessageDTO
	VotingID string            `json:"votingID"`
	Name     string            `json:"name"`
	Voter    entity.PersonInfo `json:"voter"`
}

type ChangeVoterDTO struct {
	ChatMessageDTO
	VotingID string            `json:"votingID"`
	OldName  string            `json:"oldName"`
	NewName  string            `json:"newName"`
	Voter    entity.PersonInfo `json:"voter"`
}

type TypingTextMessageDTO struct {
	ChatMessageDTO
	ChatID     string `json:"chatID"`
	SenderName string `json:"senderName"`
}
