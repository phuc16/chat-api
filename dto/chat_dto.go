package dto

import (
	"app/errors"
	"app/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type CreateChatReq struct {
	ID string `form:"id" json:"id"`
}

func (p CreateChatReq) Bind(ctx *gin.Context) (*CreateChatReq, error) {
	err := ctx.ShouldBindQuery(&p)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &p, nil
}

type GetChatActivityFromNToMReq struct {
	ID string `form:"id" json:"id"`
	X  int    `form:"x" json:"x"`
	Y  int    `form:"y" json:"y"`
}

func (p GetChatActivityFromNToMReq) Bind(ctx *gin.Context) (*GetChatActivityFromNToMReq, error) {
	err := ctx.ShouldBindQuery(&p)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &p, nil
}

type SearchByKeyWordReq struct {
	ChatID string `form:"chatID" json:"chatID"`
	Key    string `form:"key" json:"key"`
}

func (p SearchByKeyWordReq) Bind(ctx *gin.Context) (*SearchByKeyWordReq, error) {
	err := ctx.ShouldBindQuery(&p)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &p, nil
}

type GetSearchReq struct {
	ChatID    string `form:"chatID" json:"chatID"`
	MessageID string `form:"messageID" json:"messageID"`
}

func (p GetSearchReq) Bind(ctx *gin.Context) (*GetSearchReq, error) {
	err := ctx.ShouldBindQuery(&p)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &p, nil
}
