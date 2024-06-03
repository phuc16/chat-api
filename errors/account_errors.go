package errors

import "app/pkg/apperror"

const (
	CodeAccountError = 40000 + iota
	CodeAccountNotFound
	CodeAccountExists
	CodeAccountExpired
)

func AccountNotFound() *apperror.Error {
	return apperror.NewError(CodeAccountNotFound, "Tài khoản không tồn tại")
}

func AccountExists() *apperror.Error {
	return apperror.NewError(CodeAccountExists, "Tài khoản đã tồn tại")
}
