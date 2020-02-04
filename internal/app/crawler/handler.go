package crawler

import (
	"net/http"

	"github.com/dqkcode/movie-database/internal/app/types"
)

type (
	service interface {
		GetAllGenres() []string
		CrawlAllMovies() ([]*types.MovieInfo, error)
	}
	Handler struct {
		srv service
	}
)

func NewHandler(srv service) *Handler {
	return &Handler{
		srv,
	}
}

func (h *Handler) GetAllGenres(w http.ResponseWriter, r *http.Request) {
	genres := h.srv.GetAllGenres()

	types.ResponseJson(w, genres, types.Normal().Success)
}
func (h *Handler) CrawlAllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.srv.CrawlAllMovies()
	if err != nil {
		return
	}
	types.ResponseJson(w, movies, types.Normal().Success)
}
