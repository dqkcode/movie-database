package api

import (
	"github.com/dqkcode/movie-database/internal/app/policy"
	"github.com/dqkcode/movie-database/internal/app/user"
	"github.com/globalsign/mgo"
)

func NewUserService(repo *user.MongoDBRepository, policyService *policy.Service) *user.Service {

	return user.NewService(repo, policyService)
}

func NewUserServiceAndHandler(session *mgo.Session, policyService *policy.Service) (*user.Service, *user.Handler) {
	repo := user.NewMongoDBRepository(session)
	usersService := user.NewService(repo, policyService)
	return usersService, user.NewHandler(usersService)
}
