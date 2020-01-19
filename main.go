package main

import (
	"fmt"
	"github.com/dqkcode/movie-database/internal/app/auth"
	"net/http"

	"github.com/dqkcode/movie-database/internal/pkg/db/mongodb"

	"github.com/dqkcode/movie-database/internal/app/user"

	"github.com/dqkcode/movie-database/internal/pkg/http/server"

	"github.com/gorilla/mux"
)

func main() {
	ServerConf := server.LoadConfigFromEnv()
	MongoDBConf := mongodb.LoadConfigFromEnv()

	session, err := mongodb.Dial(MongoDBConf)
	if err != nil {
		panic(err)
	}
	router := mux.NewRouter()

	repo := user.NewMongoDBRepository(session)

	UserSrv := user.NewService(repo)
	UserHandler := user.NewHandler(UserSrv)

	AuthSrv := auth.NewService()
	AuthHandler := auth.NewHandler(AuthSrv)

	router.Path("/register").Methods(http.MethodPost).HandlerFunc(UserHandler.Register)
	router.Path("/login").Methods(http.MethodPost).HandlerFunc(AuthHandler.Login)

	router.Use(auth.GetUserInfoMiddleware())
	router.Path("/update").Methods(http.MethodPost).HandlerFunc(UserHandler.Update)

	router.HandleFunc("/", greet)
	server.ListenAndServe(ServerConf, router)

}

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
