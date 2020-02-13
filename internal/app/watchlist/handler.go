package watchlist

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/dqkcode/movie-database/internal/app/types"
	"github.com/gorilla/mux"
)

type (
	service interface {
		CreateWatchlist(ctx context.Context, name string) (string, error)
		AddMovieToWatchlist(ctx context.Context, movieID, watchlistID string) error
		DeleteMovieInWatchlist(ctx context.Context, movieID, watchlistID string) error
		DeleteWatchlist(ctx context.Context, watchlistID string) error
		GetWatchlistById(ctx context.Context, watchlistID string) (*WatchlistResponse, error)
		GetAllWatchlistByUserId(ctx context.Context) ([]*WatchlistResponse, error)
		ListAllMovies(ctx context.Context, watchlistID string) ([]string, error)
		UpdateStatusWatchList(ctx context.Context, watchlistID string, status bool) (string, error)
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
		MovieID string `json:"movie_id"`
	}{}
	watchlistID := mux.Vars(r)["id"]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err := h.srv.AddMovieToWatchlist(r.Context(), req.MovieID, watchlistID)
	if errors.Is(err, ErrPermissionDeny) {
		types.ResponseJson(w, "", types.Normal().PermissionDeny)
		return
	}
	if errors.Is(err, ErrWatchlistNotFound) {
		types.ResponseJson(w, "", types.Normal().NotFound)
		return
	}
	if err != nil {
		types.ResponseJson(w, "", types.Normal().Internal)
		return
	}
	types.ResponseJson(w, "", types.Normal().Success)
}

func (h *Handler) GetWatchlistById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	watchlist, err := h.srv.GetWatchlistById(r.Context(), id)

	if errors.Is(err, ErrWatchlistNotFound) {
		types.ResponseJson(w, "", types.Normal().NotFound)
		return
	}
	if errors.Is(err, ErrPermissionDeny) {
		types.ResponseJson(w, "", types.Normal().PermissionDeny)
		return

	}
	types.ResponseJson(w, watchlist, types.Normal().Success)
}

func (h *Handler) DeleteMovieInWatchlist(w http.ResponseWriter, r *http.Request) {
	watchlistID := mux.Vars(r)["watchlistID"]
	movieID := mux.Vars(r)["movieID"]

	err := h.srv.DeleteMovieInWatchlist(r.Context(), movieID, watchlistID)

	if errors.Is(err, ErrWatchlistNotFound) {
		types.ResponseJson(w, "", types.Normal().NotFound)
		return
	}
	if errors.Is(err, ErrPermissionDeny) {
		types.ResponseJson(w, "", types.Normal().PermissionDeny)
	}
	if err != nil {
		types.ResponseJson(w, "", types.Normal().Internal)
		return
	}
	types.ResponseJson(w, "", types.Normal().Success)
}

func (h *Handler) DeleteWatchlist(w http.ResponseWriter, r *http.Request) {
	watchlistID := mux.Vars(r)["id"]

	err := h.srv.DeleteWatchlist(r.Context(), watchlistID)

	if errors.Is(err, ErrWatchlistNotFound) {
		types.ResponseJson(w, "", types.Normal().NotFound)
		return
	}
	if errors.Is(err, ErrPermissionDeny) {
		types.ResponseJson(w, "", types.Normal().PermissionDeny)
	}
	if err != nil {
		types.ResponseJson(w, "", types.Normal().Internal)
		return
	}
	types.ResponseJson(w, "", types.Normal().Success)
}

func (h *Handler) GetAllMoviesInWatchlist(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	movieIDs, err := h.srv.ListAllMovies(r.Context(), id)
	if errors.Is(err, ErrPermissionDeny) {
		types.ResponseJson(w, "", types.Normal().PermissionDeny)
		return
	}
	if len(movieIDs) == 0 {
		types.ResponseJson(w, "", types.Normal().NotFound)
		return
	}
	if err != nil {
		types.ResponseJson(w, "", types.Normal().Internal)
		return
	}
	types.ResponseJson(w, movieIDs, types.Normal().Success)
}

func (h *Handler) GetAllWatchlistByUserId(w http.ResponseWriter, r *http.Request) {
	watchlists, err := h.srv.GetAllWatchlistByUserId(r.Context())
	if err != nil {
		types.ResponseJson(w, "", types.Normal().Internal)
		return
	}
	if len(watchlists) == 0 {
		types.ResponseJson(w, "", types.Normal().NotFound)
		return
	}
	types.ResponseJson(w, watchlists, types.Normal().Success)

}

func (h *Handler) UpdateStatusWatchList(w http.ResponseWriter, r *http.Request) {
	watchlistID := mux.Vars(r)["id"]
	status, err := strconv.ParseBool(r.URL.Query().Get("share"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	watchlists, err := h.srv.UpdateStatusWatchList(r.Context(), watchlistID, status)
	if errors.Is(err, ErrPermissionDeny) {
		types.ResponseJson(w, "", types.Normal().PermissionDeny)
		return
	}
	if err != nil {
		types.ResponseJson(w, "", types.Normal().Internal)
		return
	}
	if watchlists == "" {
		types.ResponseJson(w, watchlists, types.Normal().NotFound)
		return
	}
	types.ResponseJson(w, watchlists, types.Normal().Success)
}
