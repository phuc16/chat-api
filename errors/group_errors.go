package errors

import "app/pkg/apperror"

const (
	CodeGroupError = 20000 + iota
	CodeGroupNotFound
	CodeGroupExists
	CodeGroupGroupNotEnoughUser
)

func GroupNotFound() *apperror.Error {
	return apperror.NewError(CodeGroupNotFound, "group not found")
}

func GroupExists() *apperror.Error {
	return apperror.NewError(CodeGroupExists, "group exists")
}
