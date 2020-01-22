package api

import (
	"github.com/dqkcode/movie-database/internal/app/auth"
	"github.com/dqkcode/movie-database/internal/app/user"
)

func NewAuthService(usvc *user.Service) *auth.Service {
	return auth.NewService(usvc)
}

func NewAuthHandler(svc *auth.Service) *auth.Handler {
	return auth.NewHandler(svc)
}
