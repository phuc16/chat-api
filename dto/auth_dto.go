package dto

import (
	"app/errors"
	"app/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type AccountRegisterReq struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

func (r AccountRegisterReq) Bind(ctx *gin.Context) (res *AccountRegisterReq, err error) {
	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}

type AccountLoginReq struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

func (r AccountLoginReq) Bind(ctx *gin.Context) (res *AccountLoginReq, err error) {
	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}
