package dto

import (
	"app/entity"
)

type GroupMessageDTO struct {
	ID  string           `json:"id"`
	TGM TypeGroupMessage `json:"tgm"` // type group message
}

type CreateGroupDTO struct {
	GroupMessageDTO
	ChatName string              `json:"chatName"`
	Owner    entity.PersonInfo   `json:"owner"`
	Members  []entity.PersonInfo `json:"members"`
	Avatar   string              `json:"avatar"`
}

type AppendMemberGroupDTO struct {
	GroupMessageDTO
	IDChat     string `json:"idChat"`
	UserID     string `json:"userID"`
	UserName   string `json:"userName"`
	UserAvatar string `json:"userAvatar"`
}

type ChangeNameChatGroupDTO struct {
	GroupMessageDTO
	IDChat   string `json:"idChat"`
	ChatName string `json:"chatName"`
}

type ChangeAvatarGroupDTO struct {
	GroupMessageDTO
	IDChat string `json:"idChat"`
	Avatar string `json:"avatar"`
}

type UpdateSettingGroupDTO struct {
	GroupMessageDTO
	IDChat string `json:"idChat"`
	Value  bool   `json:"value"`
}

type DeleteGroupDTO struct {
	GroupMessageDTO
	IDChat string `json:"idChat"`
}
