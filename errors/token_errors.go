package errors

import "app/pkg/apperror"

const (
	CodeTokenError = 30000 + iota
	CodeTokenNotFound
	CodeTokenExists
)

func TokenNotFound() *apperror.Error {
	return apperror.NewError(CodeTokenNotFound, "Token không tồn tại")
}

func TokenExists() *apperror.Error {
	return apperror.NewError(CodeTokenExists, "Token đã tồn tại")
}
