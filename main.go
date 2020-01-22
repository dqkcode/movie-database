package main

import (
	"github.com/dqkcode/movie-database/internal/app/api"

	"github.com/dqkcode/movie-database/internal/pkg/http/server"
)

func main() {
	ServerConf := server.LoadConfigFromEnv()
	router := api.InitRouter()
	server.ListenAndServe(ServerConf, router)

}
