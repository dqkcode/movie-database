package router

import (
	"net/http"
)

type (
	Route struct {
		Path        string
		Method      string
		Handler     http.HandlerFunc
		Queries     []string
		Middlewares []Middleware
	}
	Middleware func(http.Handler) http.Handler
)
