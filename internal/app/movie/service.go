package movie

import (
	"context"
	"time"

	"github.com/dqkcode/movie-database/internal/app/types"
	"github.com/google/uuid"
)

type (
	repository interface {
		Create(ctx context.Context, movie Movie) (string, error)
		DeleteById(ctx context.Context, id string) error
		GetAllMovies(ctx context.Context) ([]*Movie, error)
		GetAllMoviesByUserId(ctx context.Context, userId string) ([]*Movie, error)
		GetMovieById(ctx context.Context, id string) (*Movie, error)
		Update(ctx context.Context, movie *Movie) error
	}

	Service struct {
		repository
	}
)

func NewService(repo repository) *Service {
	return &Service{
		repository: repo,
	}
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (string, error) {

	t, err := time.Parse("02/01/2006", req.ReleaseTime)
	if err != nil {
		//TODO CAche err
		return "", err
	}
	u := ctx.Value("user").(*types.UserInfo)

	movie := Movie{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Director:    req.Director,
		Writers:     req.Writers,
		Stars:       req.Stars,
		Casts:       req.Casts,
		MovieLength: req.MovieLength,
		ReleaseTime: t,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UserId:      u.ID,
	}
	id, err := s.repository.Create(ctx, movie)
	if err != nil {
		//TODO CAche err
		return "", err
	}
	return id, nil
}
func (s *Service) Update(ctx context.Context, id string, req UpdateRequest) error {
	u := ctx.Value("user").(*types.UserInfo)
	if u.Role != "admin" {
		m, err := s.GetMovieById(ctx, id)
		if err != nil {
			return ErrMovieNotFound
		}
		if m.UserId != u.ID {
			return ErrPermissionDeny
		}
	}

	t, err := time.Parse("02/01/2006", req.ReleaseTime)
	if err != nil {
		return err
	}
	movie := &Movie{
		ID:           id,
		Name:         req.Name,
		Rate:         req.Rate,
		Description:  req.Description,
		Director:     req.Director,
		Writers:      req.Writers,
		Stars:        req.Stars,
		TrailersPath: req.TrailersPath,
		ImagesPath:   req.ImagesPath,
		Casts:        req.Casts,
		Storyline:    req.Storyline,
		UserReviews:  req.UserReviews,
		MovieLength:  req.MovieLength,
		ReleaseTime:  t,
		UpdatedAt:    time.Now(),
	}
	if err := s.repository.Update(ctx, movie); err != nil {
		//TODO CAche err
		return err
	}
	return nil
}
func (s *Service) DeleteById(ctx context.Context, id string) error {

	err := s.repository.DeleteById(ctx, id)
	if err != nil {
		//TODO CAche err
		return err
	}
	return nil
}
func (s *Service) GetAllMovies(ctx context.Context) ([]*types.MovieInfo, error) {
	u := ctx.Value("user").(*types.UserInfo)
	if u.Role != "admin" {
		return nil, ErrPermissionDeny
	}
	movies, err := s.repository.GetAllMovies(ctx)
	if err != nil {
		//TODO CAche err
		return nil, err
	}
	var data []*types.MovieInfo
	for _, m := range movies {
		mr := m.ConvertMovieToMovieResponse()
		if err != nil {
			return nil, err
		}
		data = append(data, mr)
	}

	return data, nil
}
func (s *Service) GetAllMoviesByUserId(ctx context.Context) ([]*types.MovieInfo, error) {
	u := ctx.Value("user").(*types.UserInfo)
	movies, err := s.repository.GetAllMoviesByUserId(ctx, u.ID)
	if err != nil {
		//TODO CAche err
		return nil, err
	}
	var data []*types.MovieInfo
	for _, m := range movies {
		data = append(data, m.ConvertMovieToMovieResponse())
	}
	return data, nil
}
func (s *Service) GetMovieById(ctx context.Context, id string) (*types.MovieInfo, error) {
	movie, err := s.repository.GetMovieById(ctx, id)
	if err != nil {
		//TODO CAche err
		return nil, err
	}
	return movie.ConvertMovieToMovieResponse(), nil
}
