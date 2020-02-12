package watchlist

import (
	"context"
	"fmt"

	"github.com/dqkcode/movie-database/internal/app/types"
)

type (
	repository interface {
		CreateWatchlist(ctx context.Context, userID, name string) (string, error)
		AddMovieToWatchlist(movieID, watchlistID string) error
		DeleteMovieInWatchlist(movieID, watchlistID string) error
		GetWatchlistById(id string) (*Watchlist, error)
		ListAllMovies(id string) ([]*WatchlistMovie, error)
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
	u := ctx.Value("user").(*types.UserInfo)
	id, err := s.repo.CreateWatchlist(ctx, u.ID, name)
	if err != nil {
		return "", fmt.Errorf("failed to create watchlist %v", err)
	}
	return id, nil
}

func (s *Service) AddMovieToWatchlist(ctx context.Context, movieID, watchlistID string) error {
	if err := s.repo.AddMovieToWatchlist(movieID, watchlistID); err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteMovieInWatchlist(ctx context.Context, movieID, watchlistID string) error {
	if err := s.repo.DeleteMovieInWatchlist(movieID, watchlistID); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetWatchlistById(ctx context.Context, id string) (*WatchlistResponse, error) {

	w, err := s.repo.GetWatchlistById(id)
	if err != nil {
		return nil, err
	}
	return w.ConvertToWatchlistResponse(), nil
}

func (s *Service) ListAllMovies(ctx context.Context, id string) ([]string, error) {

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
