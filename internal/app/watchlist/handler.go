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
		CreateWatchlist(ctx context.Context, name string) (string, error)
		AddMovieToWatchlist(ctx context.Context, movieID, watchlistID string) error
		DeleteMovieInWatchlist(ctx context.Context, movieID, watchlistID string) error
		GetWatchlistById(ctx context.Context, watchlistID string) (*WatchlistResponse, error)
		ListAllMovies(ctx context.Context, watchlistID string) ([]string, error)
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

func (h *Handler) CreateWatchlist(w http.ResponseWriter, r *http.Request) {

	req := struct {
		Name string `json:"name"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	id, err := h.srv.CreateWatchlist(r.Context(), req.Name)
	if err != nil {
		types.ResponseJson(w, "", types.Normal().Internal)
	}
	types.ResponseJson(w, id, types.Normal().Success)
}

func (h *Handler) AddMovieToWatchlist(w http.ResponseWriter, r *http.Request) {
	req := struct {
		MovieID     string `json:"movie_id"`
		WatchlistID string `json:"watchlist_id"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if err := h.srv.AddMovieToWatchlist(r.Context(), req.MovieID, req.WatchlistID); err != nil {
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

func (h *Handler) DeleteMovieInWatchlist(w http.ResponseWriter, r *http.Request) {
	watchlistID := mux.Vars(r)["watchlist_id"]
	movieID := mux.Vars(r)["movie_id"]

	if err := h.srv.DeleteMovieInWatchlist(r.Context(), movieID, watchlistID); err != nil {
		if err == ErrWatchlistNotFound {
			types.ResponseJson(w, "", types.Normal().NotFound)
			return
		}
		types.ResponseJson(w, "", types.Normal().Internal)

	}
	types.ResponseJson(w, "", types.Normal().Success)
}

func (h *Handler) GetMovieInWatchlist(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	movieIDs, err := h.srv.ListAllMovies(r.Context(), id)
	if err != nil {
		if err == ErrWatchlistNotFound {
			types.ResponseJson(w, "", types.Normal().NotFound)
			return
		}
		types.ResponseJson(w, "", types.Normal().Internal)
	}
	types.ResponseJson(w, movieIDs, types.Normal().Success)
}
