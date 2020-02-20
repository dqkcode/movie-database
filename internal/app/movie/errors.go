package movie

import (
	"errors"
)

var (
	ErrMovieNotFound          = errors.New("movie not found")
	ErrPermissionDeny         = errors.New("permission deny")
	ErrUpdateMovieFailed      = errors.New("update movie failed")
	ErrInsertMovieToESFailed  = errors.New("insert movie to es failed")
	ErrIndexingDocumentFailed = errors.New("indexing movie document failed")
	ErrBadRequest             = errors.New("bad request")
)
