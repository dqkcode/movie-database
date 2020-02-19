package api

import (
	"net/http"

	"github.com/dqkcode/movie-database/internal/pkg/elasticsearch"

	"github.com/dqkcode/movie-database/internal/pkg/http/middleware"

	"github.com/dqkcode/movie-database/internal/app/auth"
	"github.com/dqkcode/movie-database/internal/pkg/db/mongodb"
	rt "github.com/dqkcode/movie-database/internal/pkg/http/router"
	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware)
	session := mongodb.InitDBSession()
	// Elasticsearch
	es := elasticsearch.NewClient()
	//Policy
	policyService := NewPolicyService()
	//User
	userService, userHandler := NewUserServiceAndHandler(session, policyService)

	//Auth
	authSrv := NewAuthService(userService)
	authHandler := auth.NewHandler(authSrv)

	//Movie
	movieService, movieHandler := NewMovieServiceAndHander(session, es)

	// Crawler
	crawlerSrv := NewCrawlerService(movieService)
	crawlerHandler := NewCrawlerHandler(crawlerSrv)

	//watchlist
	watchlistHandler := NewWatchlistHandler(session)

	//router
	routes := make([]rt.Route, 0)
	routes = append(routes, userHandler.Routes()...)
	routes = append(routes, authHandler.Routes()...)
	routes = append(routes, movieHandler.Routes()...)
	routes = append(routes, crawlerHandler.Routes()...)
	routes = append(routes, watchlistHandler.Routes()...)

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
