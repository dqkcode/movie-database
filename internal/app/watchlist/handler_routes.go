package watchlist

import (
	"net/http"

	"github.com/dqkcode/movie-database/internal/app/auth"
	"github.com/dqkcode/movie-database/internal/pkg/http/router"
)

const version1 = "/api/v1"

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			Handler:     h.AddMovieToWatchlist,
			Method:      http.MethodPost,
			Path:        version1 + "/watchlists",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
		},
		{
			Handler: h.GetWatchlistById,
			Method:  http.MethodGet,
			Path:    version1 + "/watchlists/{id}",
		},
	}
}
