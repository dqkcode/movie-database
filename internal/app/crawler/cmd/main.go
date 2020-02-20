package main

import (
	"github.com/dqkcode/movie-database/internal/app/api"
	"github.com/dqkcode/movie-database/internal/app/crawler"
)

func main() {
	s := crawler.NewServiceCLI(api.NewMovieService())

	s.CrawlAllMovies()
}
