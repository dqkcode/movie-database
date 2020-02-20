package main

import (
	"github.com/dqkcode/movie-database/internal/app/api"
	"github.com/dqkcode/movie-database/internal/pkg/http/server"
)

func main() {

	serverConf := server.LoadConfigFromEnv()
	router := api.InitRouter()
	server.ListenAndServe(serverConf, router)

	// e := movie.NewElasticsearch(elasticsearch.NewClient())
	// _, err := e.SearchMovieByName(context.Background(), "toy")
	// fmt.Println(err)

}
