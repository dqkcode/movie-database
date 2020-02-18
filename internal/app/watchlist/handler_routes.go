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
			Handler:     h.CreateWatchlist,
			Method:      http.MethodPost,
			Path:        version1 + "/watchlists",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
		},
		{
			Handler:     h.AddMovieToWatchlist,
			Method:      http.MethodPost,
			Path:        version1 + "/watchlists/{id}/movies",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
		},
		{
			Handler:     h.GetAllWatchlistByUserId,
			Method:      http.MethodGet,
			Path:        version1 + "/watchlists",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
			Queries:     []string{"ownID", `{ownID:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`},
		},
		{
			Handler:     h.GetAllMoviesInWatchlist,
			Method:      http.MethodGet,
			Path:        version1 + "/watchlists/{id}",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
			Queries:     []string{"list", "movie"},
		},
		{
			Handler:     h.GetWatchlistById,
			Method:      http.MethodGet,
			Path:        version1 + "/watchlists/{id}",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
		},

		{
			Handler:     h.DeleteWatchlist,
			Method:      http.MethodDelete,
			Path:        version1 + "/watchlists/{id}",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
		},
		{
			Handler:     h.DeleteMovieInWatchlist,
			Method:      http.MethodDelete,
			Path:        version1 + "/watchlists/{watchlistID}/{movieID}",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
		},

		{
			Handler:     h.UpdateStatusWatchList,
			Method:      http.MethodPut,
			Path:        version1 + "/watchlists/{id}",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
			Queries:     []string{"change", "status", "share", "{share:true|false}"},
		},
	}
}
