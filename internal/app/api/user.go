package api

import (
	"github.com/dqkcode/movie-database/internal/app/user"
	"github.com/globalsign/mgo"
)

func NewUserService(repo *user.MongoDBRepository) *user.Service {

	return user.NewService(repo)
}

func NewUserHandler(session *mgo.Session) *user.Handler {
	repo := user.NewMongoDBRepository(session)
	usersService := user.NewService(repo)
	return user.NewHandler(usersService)
}
