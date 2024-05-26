package dto

import (
	"app/entity"
	"app/errors"
	"app/pkg/apperror"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type AccessChatReq struct {
	UserID string `json:"userId"`
}

func (r AccessChatReq) Bind(ctx *gin.Context) (res *AccessChatReq, err error) {
	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}

func (r AccessChatReq) ToChat(ctx context.Context) (res *entity.Chat) {
	res = &entity.Chat{
		IsGroup: false,
		Users:   []string{r.UserID, entity.GetUserFromContext(ctx).ID},
	}
	return res
}

type CreateGroupReq struct {
	ChatName string   `json:"chatName"`
	Users    []string `json:"users"`
}

func (r CreateGroupReq) Bind(ctx *gin.Context) (res *CreateGroupReq, err error) {
	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}

func (r CreateGroupReq) ToChat(ctx context.Context) (res *entity.Chat) {
	res = &entity.Chat{
		ChatName:   r.ChatName,
		IsGroup:    true,
		Users:      r.Users,
		GroupAdmin: entity.GetUserFromContext(ctx).ID,
	}
	return res
}

type RenameGroupReq struct {
	ChatID   string `json:"chatId"`
	ChatName string `json:"chatName"`
}

func (r RenameGroupReq) Bind(ctx *gin.Context) (res *RenameGroupReq, err error) {
	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}

func (r RenameGroupReq) ToChat(ctx context.Context) (res *entity.Chat) {
	res = &entity.Chat{
		ID:       r.ChatID,
		ChatName: r.ChatName,
	}
	return res
}

type AddToGroupReq struct {
	UserID string `json:"userId"`
	ChatID string `json:"chatId"`
}

func (r AddToGroupReq) Bind(ctx *gin.Context) (res *AddToGroupReq, err error) {
	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}

func (r AddToGroupReq) ToChat(ctx context.Context) (res *entity.Chat) {
	res = &entity.Chat{
		ID:    r.UserID,
		Users: []string{r.UserID},
	}
	return res
}

type RemoveFromGroupReq struct {
	UserID string `json:"userId"`
	ChatID string `json:"chatId"`
}

func (r RemoveFromGroupReq) Bind(ctx *gin.Context) (res *RemoveFromGroupReq, err error) {
	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}

func (r RemoveFromGroupReq) ToChat(ctx context.Context) (res *entity.Chat) {
	res = &entity.Chat{
		ID:    r.UserID,
		Users: []string{r.UserID},
	}
	return res
}

type ChatResp struct {
	ID            string         `json:"id"`
	ChatName      string         `json:"chatName"`
	Users         []string       `json:"users"`
	IsGroup       bool           `json:"isGroup"`
	GroupAdmin    string         `json:"groupAdmin"`
	LatestMessage entity.Message `json:"latestMessage"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
}

// chuyen thanh UserIds, Users
func (r ChatResp) FromChat(e *entity.Chat) *ChatResp {
	return &ChatResp{
		ID:            e.ID,
		ChatName:      e.ChatName,
		Users:         e.Users,
		IsGroup:       e.IsGroup,
		GroupAdmin:    e.GroupAdmin,
		LatestMessage: e.LatestMessage,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

type ChatListResp struct {
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	List     []*ChatResp `json:"list"`
}
