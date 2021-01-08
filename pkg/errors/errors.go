package errors

import (
	"errors"
)

var (
	ErrUserExist = errors.New("UserExist")
	ErrUserNotFound = errors.New("UserNotFound")
	ErrForumExist = errors.New("ForumExist")
	ErrForumNotFound = errors.New("ForumNotFound")
	ErrUserOrForumNotFound = errors.New("UserOrForumNotFound")
	ErrThreadExist = errors.New("ThreadExist")
	ErrThreadNotFound = errors.New("ErrThreadNotFound")
	ErrUserOrEmailExist = errors.New("ErrUserOrEmailExist")
	ErrPostNotFound = errors.New("PostNotFound")
	ErrVoteNotFound = errors.New("VoteNotFound")
)
