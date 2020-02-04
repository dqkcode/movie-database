package api

import (
	"github.com/dqkcode/movie-database/internal/app/crawler"
	"github.com/dqkcode/movie-database/internal/app/movie"
)

func NewCrawlerService(movieSrv *movie.Service) *crawler.Service {
	return crawler.NewService(movieSrv)
}
func NewCrawlerHandler(crawlerSrv *crawler.Service) *crawler.Handler {
	return crawler.NewHandler(crawlerSrv)
}
