package crawler

import (
	"net/http"

	"github.com/dqkcode/movie-database/internal/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			Path:    "/api/v1/genres",
			Method:  http.MethodGet,
			Handler: h.GetAllGenres,
		},
		{
			Path:    "/api/v1/crawl",
			Method:  http.MethodGet,
			Handler: h.CrawlAllMovies,
		},
	}
}
