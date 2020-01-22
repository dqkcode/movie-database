package main

import (
	"net/http"

	"github.com/dqkcode/movie-database/internal/app/api"

	rt "github.com/dqkcode/movie-database/internal/pkg/http/router"

	"github.com/dqkcode/movie-database/internal/app/auth"

	"github.com/dqkcode/movie-database/internal/pkg/db/mongodb"

	"github.com/dqkcode/movie-database/internal/app/user"

	"github.com/dqkcode/movie-database/internal/pkg/http/server"

	"github.com/gorilla/mux"
)

func main() {
	ServerConf := server.LoadConfigFromEnv()
	router := mux.NewRouter()

	session := mongodb.InitDBSession()

	repo := user.NewMongoDBRepository(session)

	UserSrv := user.NewService(repo)
	UserHandler := user.NewHandler(UserSrv)

	AuthSrv := api.NewAuthService(UserSrv)
	AuthHandler := auth.NewHandler(AuthSrv)

	routes := make([]rt.Route, 0)

	routes = append(routes, UserHandler.Routes()...)
	routes = append(routes, AuthHandler.Routes()...)

	for _, r := range routes {
		h := http.Handler(r.Handler)
		for i := len(r.Middlewares) - 1; i >= 0; i-- {
			h = r.Middlewares[i](h)
		}
		router.Path(r.Path).Methods(r.Method).Handler(h).Queries(r.Queries...)

	}
	server.ListenAndServe(ServerConf, router)

}
