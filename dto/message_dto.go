package dto

import (
	"app/entity"
	"app/errors"
	"app/pkg/apperror"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type SendMessageReq struct {
	ChatID  string `json:"chatId"`
	Message string `json:"message"`
}

func (r SendMessageReq) Bind(ctx *gin.Context) (res *SendMessageReq, err error) {
	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}

func (r SendMessageReq) ToMessage(ctx context.Context) (res *entity.Message) {
	res = &entity.Message{
		Sender:  entity.GetUserFromContext(ctx).ID,
		Message: r.Message,
		ChatID:  r.ChatID,
	}
	return res
}

type MessageResp struct {
	ID        string    `bson:"id"`
	Sender    string    `bson:"sender"`
	Message   string    `bson:"message"`
	ChatID    string    `bson:"chat_id"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (r MessageResp) FromMessage(e *entity.Message) *MessageResp {
	return &MessageResp{
		ID:        r.ID,
		Sender:    r.Sender,
		Message:   r.Message,
		ChatID:    r.ChatID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

type MessageListResp struct {
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
	List     []*MessageResp `json:"list"`
}
