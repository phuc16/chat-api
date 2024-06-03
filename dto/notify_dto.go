package dto

type NotifyChat struct {
	ChatMessageDTO
	TypeNotify string
}

type NotifyGroup struct {
	GroupMessageDTO
	TypeNotify string `json:"typeNotify"`
}

type NotifyUser struct {
	UserMessageDTO
	TypeNotify string `json:"typeNotify"`
}
