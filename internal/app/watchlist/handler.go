package watchlist

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dqkcode/movie-database/internal/app/types"
	"github.com/gorilla/mux"
)

type (
	service interface {
		AddMovieToWatchlist(ctx context.Context, id string) error
		GetWatchlistById(ctx context.Context, id string) (*Watchlist, error)
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

func (h *Handler) AddMovieToWatchlist(w http.ResponseWriter, r *http.Request) {
	req := struct {
		ID string
	}{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if err := h.srv.AddMovieToWatchlist(r.Context(), req.ID); err != nil {
		types.ResponseJson(w, "", types.Normal().Internal)
	}
	types.ResponseJson(w, "", types.Normal().Success)
}

func (h *Handler) GetWatchlistById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	watchlist, err := h.srv.GetWatchlistById(r.Context(), id)
	if err != nil {
		if err == ErrWatchlistNotFound {
			types.ResponseJson(w, "", types.Normal().NotFound)
			return
		}
		types.ResponseJson(w, "", types.Normal().Internal)

	}
	types.ResponseJson(w, watchlist, types.Normal().Success)
}
