package main

import (
	"github.com/dqkcode/movie-database/internal/app/crawler"
)

func main() {
	s := crawler.NewServiceCLI()
	s.CrawlAllMovies()
}
