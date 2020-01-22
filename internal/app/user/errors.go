package user

import (
	"errors"
)

var (
	ErrUserNotFound      = errors.New("User not found")
	ErrDBQuery           = errors.New("DB Query error")
	ErrUpdateUserFailed  = errors.New("Update user failed")
	ErrUserAlreadyExist  = errors.New("User already exist")
	ErrGenPasswordFailed = errors.New("Generate password failed")
	ErrCreateUserFailed  = errors.New("Create user failed")
	ErrPermissionDeny    = errors.New("Permission deny")
	Err                  = errors.New("")
)
