package watchlist

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/dqkcode/movie-database/internal/app/types"
)

type (
	repository interface {
		CreateWatchlist(ctx context.Context, userID, name string) (string, error)
		AddMovieToWatchlist(movieID, watchlistID string) error
		DeleteMovieInWatchlist(movieID, watchlistID string) error
		DeleteWatchlist(watchlistID string) error
		GetWatchlistById(id string) (*Watchlist, error)
		ListAllMovies(id string) ([]*WatchlistMovie, error)
		GetAllWatchlistByUserId(id string) ([]*Watchlist, error)
		UpdateStatusWatchList(ctx context.Context, watchlistID string, status bool) error
	}
	Service struct {
		repo repository
	}
)

func NewService(repository repository) *Service {
	return &Service{
		repo: repository,
	}
}

func (s *Service) CreateWatchlist(ctx context.Context, name string) (string, error) {
	u := ctx.Value(types.UserContextKey).(*types.UserInfo)
	id, err := s.repo.CreateWatchlist(ctx, u.ID, name)
	if err != nil {
		return "", fmt.Errorf("failed to create watchlist %w", err)
	}
	return id, nil
}

func (s *Service) AddMovieToWatchlist(ctx context.Context, movieID, watchlistID string) error {
	u := ctx.Value(types.UserContextKey).(*types.UserInfo)
	w, err := s.repo.GetWatchlistById(watchlistID)
	if errors.Is(err, ErrWatchlistNotFound) {
		return err
	}
	if err != nil {
		return fmt.Errorf("err from GetWatchlistById :%w", err)
	}
	if w.UserID != u.ID && u.Role != "admin" {
		return ErrPermissionDeny
	}
	if err := s.repo.AddMovieToWatchlist(movieID, watchlistID); err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteMovieInWatchlist(ctx context.Context, movieID, watchlistID string) error {
	u := ctx.Value(types.UserContextKey).(*types.UserInfo)
	w, err := s.repo.GetWatchlistById(watchlistID)
	if err != nil {
		return fmt.Errorf("err from GetWatchlistById: err =  %w", err)
	}
	if u.ID != w.UserID && u.Role != "admin" {
		return ErrPermissionDeny
	}
	if err := s.repo.DeleteMovieInWatchlist(movieID, watchlistID); err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteWatchlist(ctx context.Context, watchlistID string) error {
	u := ctx.Value(types.UserContextKey).(*types.UserInfo)
	w, err := s.repo.GetWatchlistById(watchlistID)
	if err != nil {
		return fmt.Errorf("err from GetWatchlistById err =  %w", err)
	}
	if u.ID != w.UserID && u.Role != "admin" {
		return ErrPermissionDeny
	}
	if err := s.repo.DeleteWatchlist(watchlistID); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetWatchlistById(ctx context.Context, id string) (*WatchlistResponse, error) {
	u := ctx.Value(types.UserContextKey).(*types.UserInfo)
	w, err := s.repo.GetWatchlistById(id)
	if err != nil {
		return nil, err
	}
	if u.ID != w.UserID && w.Share == false && u.Role != "admin" {
		return nil, ErrPermissionDeny
	}
	return w.ConvertToWatchlistResponse(), nil
}

func (s *Service) ListAllMovies(ctx context.Context, id string) ([]string, error) {
	w, err := s.repo.GetWatchlistById(id)
	if err != nil {
		return nil, err
	}
	var user *types.UserInfo
	if u := ctx.Value(types.UserContextKey); u != nil || reflect.ValueOf(u).IsNil() == false {
		user = u.(*types.UserInfo)

	}
	if w.Share == false && w.UserID != user.ID && user.Role != "admin" {
		return nil, ErrPermissionDeny
	}
	list, err := s.repo.ListAllMovies(id)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, m := range list {
		ids = append(ids, m.MovieID)
	}
	return ids, nil
}

func (s *Service) GetAllWatchlistByUserId(ctx context.Context, userID string) ([]*WatchlistResponse, error) {
	u := ctx.Value(types.UserContextKey).(*types.UserInfo)
	if u.ID != userID && u.Role != "admin" {
		return nil, ErrPermissionDeny
	}
	list, err := s.repo.GetAllWatchlistByUserId(userID)
	if err != nil {
		return nil, err
	}
	var result []*WatchlistResponse
	for _, m := range list {
		result = append(result, m.ConvertToWatchlistResponse())
	}
	return result, nil
}

func (s *Service) UpdateStatusWatchList(ctx context.Context, watchlistID string, status bool) (string, error) {
	u := ctx.Value(types.UserContextKey).(*types.UserInfo)
	w, err := s.repo.GetWatchlistById(watchlistID)
	if err != nil {
		return "", err
	}
	if w.UserID != u.ID && u.Role != "admin" {
		return "", ErrPermissionDeny
	}
	if err := s.repo.UpdateStatusWatchList(ctx, watchlistID, status); err != nil {
		return "", err
	}
	return watchlistID, nil
}
