package errors

import "app/pkg/apperror"

const (
	CodeGroupError = 20000 + iota
	CodeGroupNotFound
	CodeGroupExists
	CodeGroupGroupNotEnoughUser
)

func GroupNotFound() *apperror.Error {
	return apperror.NewError(CodeGroupNotFound, "Nhóm không tồn tại")
}

func GroupExists() *apperror.Error {
	return apperror.NewError(CodeGroupExists, "Nhóm đã tồn tại")
}
