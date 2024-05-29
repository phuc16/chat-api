package errors

import "app/pkg/apperror"

const (
	CodeVotingError = 20000 + iota
	CodeVotingNotFound
	CodeVotingExists
	CodeVotingVotingNotEnoughUser
)

func VotingNotFound() *apperror.Error {
	return apperror.NewError(CodeVotingNotFound, "voting not found")
}

func VotingExists() *apperror.Error {
	return apperror.NewError(CodeVotingExists, "voting exists")
}
