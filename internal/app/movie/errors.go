package movie

import (
	"errors"
)

var (
	ErrMovieNotFound     = errors.New("movie not found")
	ErrPermissionDeny    = errors.New("permission deny")
	ErrUpdateMovieFailed = errors.New("update movie failed")
)
