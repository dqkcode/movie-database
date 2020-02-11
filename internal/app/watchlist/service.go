package watchlist

import (
	"context"

	"github.com/dqkcode/movie-database/internal/app/types"
)

type (
	repository interface {
		AddMovieToWatchlist(userID, movieID string) error
		GetWatchlistById(id string) (*Watchlist, error)
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

func (s *Service) AddMovieToWatchlist(ctx context.Context, movieID string) error {
	u := ctx.Value("user").(*types.UserInfo)
	if err := s.repo.AddMovieToWatchlist(u.ID, movieID); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetWatchlistById(ctx context.Context, id string) (*Watchlist, error) {

	w, err := s.repo.GetWatchlistById(id)
	if err != nil {
		return nil, err
	}
	return w, nil
}
