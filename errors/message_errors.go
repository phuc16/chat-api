package errors

import "app/pkg/apperror"

const (
	CodeMessageError = 20000 + iota
	CodeMessageNotFound
	CodeMessageExists
	CodeMessageGroupNotEnoughUser
)

func MessageNotFound() *apperror.Error {
	return apperror.NewError(CodeMessageNotFound, "message not found")
}

func MessageExists() *apperror.Error {
	return apperror.NewError(CodeMessageExists, "message exists")
}
