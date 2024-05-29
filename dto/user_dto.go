package dto

import (
	"app/errors"
	"app/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type CreateUserReq struct {
	ID string `form:"id" json:"id"`
}

func (r CreateUserReq) Bind(ctx *gin.Context) (res *CreateUserReq, err error) {
	err = ctx.ShouldBindQuery(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}

type UpdateAvatarAsyncReq struct {
	OldAvatar string `form:"oldAvatar" json:"oldAvatar"`
	NewAvatar string `form:"newAvatar" json:"newAvatar"`
}

func (p UpdateAvatarAsyncReq) Bind(ctx *gin.Context) (*UpdateAvatarAsyncReq, error) {
	err := ctx.ShouldBindQuery(&p)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &p, nil
}
