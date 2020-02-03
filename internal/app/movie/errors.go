package movie

import (
	"errors"
)

var (
	ErrMovieNotFound     = errors.New("Movie not found")
	ErrPermissionDeny    = errors.New("Permission deny")
	ErrUpdateMovieFailed = errors.New("Update movie failed")
)
