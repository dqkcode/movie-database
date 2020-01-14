package main

import (
	"fmt"
	"net/http"

	"github.com/dqkcode/movie-database/internal/pkg/http/server"
	"github.com/gorilla/mux"
)

func main() {
	conf := server.LoadConfigFromEnv()

	router := mux.NewRouter()
	router.HandleFunc("/", greet)
	server.ListenAndServe(conf, router)

}

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
