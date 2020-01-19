package api

import (
	"github.com/dqkcode/movie-database/internal/app/user"
	"github.com/dqkcode/movie-database/internal/pkg/db/mongodb"
)

func NewUserServive() *user.Service {

	conf := mongodb.LoadConfigFromEnv()
	session, err := mongodb.Dial(conf)
	if err != nil {
		panic(err)
	}

	repo := user.NewMongoDBRepository(session)
	srv := user.NewService(repo)

	return srv
}
