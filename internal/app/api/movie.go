package api

import (
	"github.com/dqkcode/movie-database/internal/app/movie"
	"github.com/globalsign/mgo"
)

func NewMovieService(repo *movie.MongoDBRepository) *movie.Service {
	return movie.NewService(repo)
}
func NewMovieHander(session *mgo.Session) *movie.Handler {
	repo := movie.NewMongoDBRepository(session)
	srv := movie.NewService(repo)
	return movie.NewHandler(srv)
}
