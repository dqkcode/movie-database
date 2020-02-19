package main

import (
	"fmt"

	"github.com/dqkcode/movie-database/internal/app/types"

	"github.com/dqkcode/movie-database/internal/app/movie"
	"github.com/dqkcode/movie-database/internal/pkg/elasticsearch"
)

func main() {

	// serverConf := server.LoadConfigFromEnv()
	// router := api.InitRouter()
	// server.ListenAndServe(serverConf, router)

	e := movie.NewElasticsearch(elasticsearch.NewClient())
	m := &types.MovieInfo{
		ID:   "ba3a97a8-5309-11ea-8d77-2e728ce88125",
		Name: "name 1",
	}
	err := e.InsertMovies(m)
	fmt.Println(err)

}
