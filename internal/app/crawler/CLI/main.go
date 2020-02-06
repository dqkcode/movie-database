package main

import (
	"context"

	"github.com/dqkcode/movie-database/internal/app/api"
	"github.com/dqkcode/movie-database/internal/app/crawler"
)

func main() {
	s := crawler.NewServiceCLI(api.NewMovieServiceWithMongoRepo())
	newCtx := context.Background()
	s.CrawlAllMovies(newCtx)
}
