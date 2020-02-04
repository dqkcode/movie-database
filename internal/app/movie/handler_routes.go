package movie

import (
	"net/http"

	"github.com/dqkcode/movie-database/internal/app/auth"

	"github.com/dqkcode/movie-database/internal/pkg/http/router"
)

const version1 = "/api/v1"

func (h *Handler) Routes() []router.Route {

	return []router.Route{
		{
			Handler:     h.CreateMovie,
			Method:      http.MethodPost,
			Path:        version1 + "/movies",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
		},
		{
			Handler:     h.DeleteMovieById,
			Method:      http.MethodDelete,
			Path:        version1 + "/movies/{id}",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
		},
		{
			Handler:     h.Update,
			Method:      http.MethodPut,
			Path:        version1 + "/movies/{id}",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
		},
		{
			Handler:     h.GetAllMovies,
			Method:      http.MethodGet,
			Path:        version1 + "/movies",
			Queries:     []string{"role", "admin"},
			Middlewares: []router.Middleware{auth.AuthMiddleware},
		},
		{
			Handler:     h.GetAllMoviesByUserId,
			Method:      http.MethodGet,
			Path:        version1 + "/movies",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
		},
		{
			Handler:     h.GetMovieById,
			Method:      http.MethodGet,
			Path:        version1 + "/movies/{id}",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
		},
	}
}
