package main

import (
	"fmt"
	"net/http"

	"github.com/globalsign/mgo"

	"github.com/dqkcode/movie-database/internal/pkg/db/mongodb"

	"github.com/dqkcode/movie-database/internal/app/user"

	"github.com/dqkcode/movie-database/internal/pkg/http/server"
	"github.com/gorilla/mux"
)

func main() {
	ServerConf := server.LoadConfigFromEnv()
	MongoDBConf := mongodb.LoadConfigFromEnv()

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    MongoDBConf.Addrs,
		Database: MongoDBConf.Database,
		Password: MongoDBConf.Password,
		Username: MongoDBConf.Username,
		Timeout:  MongoDBConf.Timeout,
	})
	if err != nil {
		panic(err)
	}
	router := mux.NewRouter()

	repo := user.NewMongoDBRepository(session)
	srv := user.NewService(repo)
	UserHandler := user.NewHandler(srv)

	router.Path("/register").Methods(http.MethodPost).HandlerFunc(UserHandler.Register)

	router.HandleFunc("/", greet)
	server.ListenAndServe(ServerConf, router)

}

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
