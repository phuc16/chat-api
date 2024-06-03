package dto

import (
	"app/errors"
	"app/pkg/apperror"
	"time"

	"github.com/gin-gonic/gin"
)

type AccountCheckPhoneNumberReq struct {
	PhoneNumber string `json:"phoneNumber"`
}

func (r AccountCheckPhoneNumberReq) Bind(ctx *gin.Context) (res *AccountCheckPhoneNumberReq, err error) {
	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}

type AccountResetPasswordReq struct {
	PhoneNumber string `json:"phoneNumber"`
	NewPassword string `json:"newPassword"`
}

func (r AccountResetPasswordReq) Bind(ctx *gin.Context) (res *AccountResetPasswordReq, err error) {
	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}

type AccountChangePasswordReq struct {
	CurPassword string `json:"curPassword"`
	NewPassword string `json:"newPassword"`
}

func (r AccountChangePasswordReq) Bind(ctx *gin.Context) (res *AccountChangePasswordReq, err error) {
	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}

type AccountChangeAvatarReq struct {
	NewAvatar string `json:"newAvatar"`
}

func (r AccountChangeAvatarReq) Bind(ctx *gin.Context) (res *AccountChangeAvatarReq, err error) {
	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}

type AccountChangeProfileReq struct {
	UserName string    `json:"userName"`
	Gender   bool      `json:"gender"`
	Birthday time.Time `json:"birthday"`
}

func (r AccountChangeProfileReq) Bind(ctx *gin.Context) (res *AccountChangeProfileReq, err error) {
	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}
