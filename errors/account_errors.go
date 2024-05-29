package errors

import "app/pkg/apperror"

const (
	CodeAccountError = 40000 + iota
	CodeAccountNotFound
	CodeAccountExists
	CodeAccountExpired
)

func AccountNotFound() *apperror.Error {
	return apperror.NewError(CodeAccountNotFound, "account not found")
}

func AccountExists() *apperror.Error {
	return apperror.NewError(CodeAccountExists, "account exists")
}