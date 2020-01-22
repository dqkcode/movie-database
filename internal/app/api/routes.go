package api

import (
	"net/http"

	"github.com/dqkcode/movie-database/internal/pkg/http/middleware"

	"github.com/dqkcode/movie-database/internal/app/auth"
	"github.com/dqkcode/movie-database/internal/app/user"
	"github.com/dqkcode/movie-database/internal/pkg/db/mongodb"
	rt "github.com/dqkcode/movie-database/internal/pkg/http/router"
	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware)
	session := mongodb.InitDBSession()
	repo := user.NewMongoDBRepository(session)

	//User
	usersService := user.NewService(repo)
	usersHandler := user.NewHandler(usersService)

	//Auth
	authSrv := NewAuthService(usersService)
	authHandler := auth.NewHandler(authSrv)

	routes := make([]rt.Route, 0)
	routes = append(routes, usersHandler.Routes()...)
	routes = append(routes, authHandler.Routes()...)

	//add attributes to router
	for _, r := range routes {
		h := http.Handler(r.Handler)
		for i := len(r.Middlewares) - 1; i >= 0; i-- {
			h = r.Middlewares[i](h)
		}
		router.Path(r.Path).Methods(r.Method).Handler(h).Queries(r.Queries...)

	}
	return router
}
