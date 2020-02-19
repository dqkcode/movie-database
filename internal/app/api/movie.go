package api

import (
	"github.com/dqkcode/movie-database/internal/app/movie"
	"github.com/elastic/go-elasticsearch/v7"

	"github.com/globalsign/mgo"
)

func NewMovieServiceAndHander(session *mgo.Session, c *elasticsearch.Client) (*movie.Service, *movie.Handler) {
	repo := movie.NewMongoDBRepository(session)
	es := movie.NewElasticsearch(c)
	s := movie.NewService(repo, es)
	return s, movie.NewHandler(s)
}
