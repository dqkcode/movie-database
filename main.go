package main

import (
	"fmt"
	"net/http"

	"github.com/dqkcode/movie-database/internal/app/auth"
	"github.com/sirupsen/logrus"

	"github.com/dqkcode/movie-database/internal/pkg/db/mongodb"

	"github.com/dqkcode/movie-database/internal/app/user"

	"github.com/dqkcode/movie-database/internal/pkg/http/server"

	"github.com/gorilla/mux"
)

func main() {
	ServerConf := server.LoadConfigFromEnv()
	MongoDBConf := mongodb.LoadConfigFromEnv()

	logrus.Infof("MongoDb env: %v", MongoDBConf)
	// fmt.Printf("MongoDb env: %v", MongoDBConf)
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

	router.Handle("/update", auth.UserInfoMiddleware(http.HandlerFunc(UserHandler.Update))).Methods(http.MethodPost)

	router.HandleFunc("/", greet)
	server.ListenAndServe(ServerConf, router)

}

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
