package movie

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/dqkcode/movie-database/internal/app/types"
)

type (
	service interface {
		Create(ctx context.Context, req CreateRequest) (string, error)
		DeleteById(ctx context.Context, id string) error
		GetAllMovies(ctx context.Context, req FindRequest) ([]*types.MovieInfo, error)
		GetMovieById(ctx context.Context, id string) (*types.MovieInfo, error)
		Update(ctx context.Context, id string, movie UpdateRequest) error
		SearchMovieByName(ctx context.Context, movieName string) ([]types.MovieInfo, error)
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
		types.ResponseJson(w, "", types.Normal().Internal)
		return
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
		return
	}
	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.srv.Update(r.Context(), id, req)
	if errors.Is(err, ErrMovieNotFound) {
		types.ResponseJson(w, "", types.Movie().NotFound)
		return
	}
	if errors.Is(err, ErrPermissionDeny) {
		types.ResponseJson(w, "", types.Normal().PermissionDeny)
		return
	}
	if err != nil {
		types.ResponseJson(w, "", types.Normal().Internal)
		return
	}
	types.ResponseJson(w, "", types.Normal().Success)
}

func (h *Handler) DeleteMovieById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		types.ResponseJson(w, "", types.Normal().BadRequest)
		return
	}
	err := h.srv.DeleteById(r.Context(), id)
	if errors.Is(err, ErrMovieNotFound) {
		types.ResponseJson(w, "", types.Movie().NotFound)
		return
	}
	if errors.Is(err, ErrPermissionDeny) {
		types.ResponseJson(w, "", types.Normal().PermissionDeny)
		return
	}
	if err != nil {
		types.ResponseJson(w, "", types.Movie().DeleteFailed)
		return
	}
	types.ResponseJson(w, "", types.Normal().Success)
}

func (h *Handler) GetAllMovies(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	movieLength, _ := strconv.Atoi(queries.Get("max_length"))
	offset, _ := strconv.Atoi(queries.Get("offset"))
	limit, _ := strconv.Atoi(queries.Get("limit"))
	rate, _ := strconv.ParseFloat(queries.Get("rate"), 8)
	req := FindRequest{
		Rate:        rate,
		Offset:      offset,
		Limit:       limit,
		MovieLength: movieLength,
		Name:        queries.Get("name"),
		ReleaseTime: queries.Get("release_time"),
		CreatedByID: queries.Get("create_by_id"),
		Directors:   queries["directors"],
		Casts:       queries["casts"],
		Writers:     queries["writers"],
		Genres:      queries["genres"],
		Selects:     queries["selects"],
		SortBy:      queries["sort_by"],
	}
	movies, err := h.srv.GetAllMovies(r.Context(), req)
	if errors.Is(err, ErrMovieNotFound) {
		types.ResponseJson(w, "", types.Movie().NotFound)
		return
	}
	if err != nil {
		types.ResponseJson(w, "", types.Normal().NotFound)
		return
	}
	types.ResponseJson(w, movies, types.Normal().Success)
}

func (h *Handler) GetMovieById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		types.ResponseJson(w, "", types.Normal().BadRequest)
		return
	}
	movie, err := h.srv.GetMovieById(r.Context(), id)
	if errors.Is(err, ErrMovieNotFound) {
		types.ResponseJson(w, "", types.Movie().NotFound)
		return
	}
	if err != nil {
		types.ResponseJson(w, "", types.Normal().Internal)
		return
	}
	types.ResponseJson(w, movie, types.Normal().Success)
}

func (h *Handler) SearchMovieByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		types.ResponseJson(w, "", types.Normal().BadRequest)
		return
	}
	movies, err := h.srv.SearchMovieByName(r.Context(), name)
	if errors.Is(err, ErrMovieNotFound) {
		types.ResponseJson(w, "", types.Movie().NotFound)
		return
	}
	if err != nil {
		types.ResponseJson(w, "", types.Normal().Internal)
		return
	}
	types.ResponseJson(w, movies, types.Normal().Success)
}
