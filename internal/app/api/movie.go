package api

import (
	"github.com/dqkcode/movie-database/internal/app/movie"
	"github.com/globalsign/mgo"
)

// func NewMovieService(repo *movie.MongoDBRepository) *movie.Service {
// 	return movie.NewService(repo)
// }
func NewMovieService(session *mgo.Session) *movie.Service {
	repo := movie.NewMongoDBRepository(session)
	return movie.NewService(repo)
}
func NewMovieHander(srv *movie.Service) *movie.Handler {

	return movie.NewHandler(srv)
}
