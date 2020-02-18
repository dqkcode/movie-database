package user

import (
	"errors"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrDB                = errors.New("DB error")
	ErrUpdateUserFailed  = errors.New("update user failed")
	ErrUserAlreadyExist  = errors.New("user already exist")
	ErrGenPasswordFailed = errors.New("generate password failed")
	ErrCreateUserFailed  = errors.New("create user failed")
	ErrPermissionDeny    = errors.New("permission deny")
)
