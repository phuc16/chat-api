package errors

import "app/pkg/apperror"

const (
	CodeMessageError = 20000 + iota
	CodeMessageNotFound
	CodeMessageExists
	CodeMessageGroupNotEnoughUser
)

func MessageNotFound() *apperror.Error {
	return apperror.NewError(CodeMessageNotFound, "Tin nhắn không tồn tại")
}

func MessageExists() *apperror.Error {
	return apperror.NewError(CodeMessageExists, "Tin nhắn đã tồn tại")
}
