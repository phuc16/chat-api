package errors

import "app/pkg/apperror"

const (
	CodeChatError = 20000 + iota
	CodeChatNotFound
	CodeChatExists
	CodeChatGroupNotEnoughUser
)

func ChatNotFound() *apperror.Error {
	return apperror.NewError(CodeChatNotFound, "chat not found")
}

func ChatExists() *apperror.Error {
	return apperror.NewError(CodeChatExists, "chat exists")
}

func ChatGroupNotEnoughUser() *apperror.Error {
	return apperror.NewError(CodeChatGroupNotEnoughUser, "group should contain more than 2 users")
}
