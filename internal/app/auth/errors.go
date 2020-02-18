package auth

import (
	"errors"
)

var (
	ErrUserIsLocked = errors.New("user was locked")
	ErrCompareToken = errors.New("can not compare token")
	ErrTokenInvalid = errors.New("the token is invalid")
)
