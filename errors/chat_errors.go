package errors

import "app/pkg/apperror"

const (
	CodeChatError = 20000 + iota
	CodeChatNotFound
	CodeChatExists
	CodeChatGroupNotEnoughUser
)

func ChatNotFound() *apperror.Error {
	return apperror.NewError(CodeChatNotFound, "Hội thoại không tồn tại")
}

func ChatExists() *apperror.Error {
	return apperror.NewError(CodeChatExists, "Hội thoại đã tồn tại")
}

func ChatGroupNotEnoughUser() *apperror.Error {
	return apperror.NewError(CodeChatGroupNotEnoughUser, "group should contain more than 2 users")
}
