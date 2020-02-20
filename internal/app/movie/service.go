package movie

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/dqkcode/movie-database/internal/app/types"

	"github.com/google/uuid"
)

type (
	repository interface {
		Create(ctx context.Context, movie Movie) (string, error)
		DeleteById(ctx context.Context, id string) error
		GetAllMovies(ctx context.Context, req FindRequest) ([]*Movie, error)
		GetMovieByName(ctx context.Context, name string) (*Movie, error)
		GetMovieById(ctx context.Context, id string) (*Movie, error)
		Update(ctx context.Context, movie *Movie) error
	}
	searchEngine interface {
		InsertMovies(ctx context.Context, movie *types.MovieInfo) error
		SearchMovieByName(ctx context.Context, movieName string) ([]types.MovieInfo, error)
	}

	Service struct {
		repository   repository
		searchEngine searchEngine
	}
)

func NewService(repo repository, s searchEngine) *Service {
	return &Service{
		repository:   repo,
		searchEngine: s,
	}
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (string, error) {
	u := ctx.Value(types.UserContextKey).(*types.UserInfo)
	t := time.Now()
	movie := Movie{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Rate:        req.Rate,
		Directors:   req.Directors,
		Writers:     req.Writers,
		Casts:       req.Casts,
		MovieLength: req.MovieLength,
		ReleaseTime: req.ReleaseTime,
		CreatedAt:   &t,
		UpdatedAt:   &t,
		UserId:      u.ID,
	}
	id, err := s.repository.Create(ctx, movie)
	if err != nil {
		return "", fmt.Errorf("failed to create movie, err: %w", err)
	}
	if err := s.searchEngine.InsertMovies(ctx, movie.ConvertMovieToMovieResponse()); err != nil {
		logrus.Errorf("failed to insert movie to es, err: %w", err)
	}
	return id, nil
}

//CreateMovie for crawler use
func (s *Service) CreateMovie(m types.MovieInfo) error {
	t := time.Now()
	movie := Movie{
		ID:           uuid.New().String(),
		Name:         m.Name,
		MovieLength:  m.MovieLength,
		ReleaseTime:  m.ReleaseTime,
		Directors:    m.Directors,
		Writers:      m.Writers,
		Rate:         m.Rate,
		Genres:       m.Genres,
		Casts:        m.Casts,
		Storyline:    m.Storyline,
		ImagesPath:   m.ImagesPath,
		TrailersPath: m.TrailersPath,
		CreatedAt:    &t,
		UpdatedAt:    &t,
	}
	newCtx := context.Background()
	_, err := s.repository.Create(newCtx, movie)
	if err != nil {
		return fmt.Errorf("failed to create movie, err: %w", err)
	}
	if err := s.searchEngine.InsertMovies(newCtx, movie.ConvertMovieToMovieResponse()); err != nil {
		logrus.Errorf("failed to insert movie to es, err: %w", err)
	}
	return nil
}

func (s *Service) Update(ctx context.Context, id string, req UpdateRequest) error {
	u := ctx.Value(types.UserContextKey).(*types.UserInfo)
	m, err := s.GetMovieById(ctx, id)
	if err != nil {
		return err
	}
	if u.Role != "admin" && m.UserId != u.ID {
		return ErrPermissionDeny
	}
	t := time.Now()
	movie := &Movie{
		ID:           id,
		Name:         req.Name,
		Rate:         req.Rate,
		Directors:    req.Directors,
		Writers:      req.Writers,
		TrailersPath: req.TrailersPath,
		ImagesPath:   req.ImagesPath,
		Casts:        req.Casts,
		Storyline:    req.Storyline,
		UserReviews:  req.UserReviews,
		MovieLength:  req.MovieLength,
		ReleaseTime:  req.ReleaseTime,
		UpdatedAt:    &t,
	}
	if err := s.repository.Update(ctx, movie); err != nil {
		return fmt.Errorf("failed to update movie, err: %w", err)
	}
	return nil
}

func (s *Service) DeleteById(ctx context.Context, id string) error {
	u := ctx.Value(types.UserContextKey).(*types.UserInfo)
	m, err := s.GetMovieById(ctx, id)
	if err != nil {
		return err
	}
	if u.Role != "admin" && m.UserId != u.ID {
		return ErrPermissionDeny
	}
	err = s.repository.DeleteById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete movie, err: %w", err)
	}
	return nil
}

func (s *Service) GetAllMovies(ctx context.Context, req FindRequest) ([]*types.MovieInfo, error) {
	if req.Offset > 50 {
		req.Offset = 50
	}
	movies, err := s.repository.GetAllMovies(ctx, req)
	if err == ErrMovieNotFound {
		return nil, ErrMovieNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get all movie, err: %w", err)
	}
	var data []*types.MovieInfo
	for _, m := range movies {
		mr := m.ConvertMovieToMovieResponse()
		if err != nil {
			return nil, fmt.Errorf("failed to convert movie, err: %w", err)
		}
		data = append(data, mr)
	}
	return data, nil
}

func (s *Service) GetMovieByName(name string) (*types.MovieInfo, error) {
	ctx := context.Background()
	movie, err := s.repository.GetMovieByName(ctx, name)
	if errors.Is(err, ErrMovieNotFound) {
		return nil, types.ErrMovieNotFound
	}
	if err != nil {
		return nil, err
	}
	return movie.ConvertMovieToMovieResponse(), nil
}

func (s *Service) GetMovieById(ctx context.Context, id string) (*types.MovieInfo, error) {
	movie, err := s.repository.GetMovieById(ctx, id)
	if err != nil {
		return nil, err
	}
	return movie.ConvertMovieToMovieResponse(), nil
}

func (s *Service) SearchMovieByName(ctx context.Context, name string) ([]types.MovieInfo, error) {
	movies, err := s.searchEngine.SearchMovieByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return movies, nil
}
