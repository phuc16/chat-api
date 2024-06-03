package dto

import (
	"app/entity"
	"app/errors"
	"app/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type CreateGroupReq struct {
	entity.Group
}

func (r CreateGroupReq) Bind(ctx *gin.Context) (res *CreateGroupReq, err error) {
	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}

type GetGroupInfoReq struct {
	GroupID string `form:"groupID" json:"groupID"`
}

func (p GetGroupInfoReq) Bind(ctx *gin.Context) (*GetGroupInfoReq, error) {
	err := ctx.ShouldBindQuery(&p)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &p, nil
}
