package movie

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/dqkcode/movie-database/internal/app/types"
)

type (
	service interface {
		Create(ctx context.Context, req CreateRequest) (string, error)
		DeleteById(ctx context.Context, id string) error
		GetAllMovies(ctx context.Context) ([]*types.MovieInfo, error)
		GetAllMoviesByUserId(ctx context.Context) (*[]Movie, error)
		Update(ctx context.Context, id string, movie UpdateRequest) error
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

func (h *Handler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := h.srv.Create(r.Context(), req)
	if err != nil {
		types.ResponseJson(w, "", types.Normal().BadRequest)

	}
	data := map[string]string{
		"id": id,
	}
	types.ResponseJson(w, data, types.Normal().Success)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		types.ResponseJson(w, "", types.Normal().BadRequest)
	}
	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.srv.Update(r.Context(), id, req); err != nil {
		types.ResponseJson(w, "", types.Normal().Internal)
		return

	}
	types.ResponseJson(w, "", types.Normal().Success)
	return
}

func (h *Handler) DeleteMovieById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		types.ResponseJson(w, "", types.Normal().BadRequest)
		return
	}
	err := h.srv.DeleteById(r.Context(), id)
	if err != nil {
		types.ResponseJson(w, "", types.Movie().DeleteFailed)
		return
	}
	types.ResponseJson(w, "", types.Normal().Success)
	return
}

func (h *Handler) GetAllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.srv.GetAllMovies(r.Context())
	if err != nil {
		if err == ErrPermissionDeny {
			types.ResponseJson(w, "", types.Normal().PermissionDeny)
			return
		}
		types.ResponseJson(w, "", types.Normal().NotFound)
		return
	}
	types.ResponseJson(w, movies, types.Normal().Success)
	return
}
func (h *Handler) GetAllMoviesByUserId(w http.ResponseWriter, r *http.Request) {
	movies, err := h.srv.GetAllMoviesByUserId(r.Context())
	if err != nil {
		if err == ErrPermissionDeny {
			types.ResponseJson(w, "", types.Normal().PermissionDeny)
			return
		}
		types.ResponseJson(w, "", types.Normal().NotFound)
		return
	}
	types.ResponseJson(w, movies, types.Normal().Success)
	return
}
