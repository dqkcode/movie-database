package api

import (
	"github.com/dqkcode/movie-database/internal/app/user"
)

func NewUserServive() {
	repo := user.NewMongoDBRepository(session)
	srv := user.NewService(repo)
}
