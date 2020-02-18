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
			// Create
			Handler:     h.CreateMovie,
			Method:      http.MethodPost,
			Path:        version1 + "/movies",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
		},
		{
			//Get All
			Handler: h.GetAllMovies,
			Method:  http.MethodGet,
			Path:    version1 + "/movies",
		},
		{
			// Get By ID
			Handler: h.GetMovieById,
			Method:  http.MethodGet,
			Path:    version1 + "/movies/{id}",
		},
		{
			// Update By ID
			Handler:     h.Update,
			Method:      http.MethodPut,
			Path:        version1 + "/movies/{id}",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
		},
		{
			// Delete By ID
			Handler:     h.DeleteMovieById,
			Method:      http.MethodDelete,
			Path:        version1 + "/movies/{id}",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
		},
	}
}
