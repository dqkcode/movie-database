package movie

import (
	"context"
	"strconv"
	"time"

	"github.com/dqkcode/movie-database/internal/app/types"
	"github.com/google/uuid"
)

type (
	repository interface {
		Create(ctx context.Context, movie Movie) (string, error)
		DeleteById(ctx context.Context, id string) error
		GetAllMovies(ctx context.Context, req FindRequest) ([]*Movie, error)
		GetAllMoviesByUserId(ctx context.Context, userId string) ([]*Movie, error)
		GetMovieByName(ctx context.Context, name string) (*Movie, error)
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

	// t, err := time.Parse("02/01/2006", req.ReleaseTime)
	// if err != nil {
	// 	//TODO CAche err
	// 	return "", err
	// }
	u := ctx.Value("user").(*types.UserInfo)
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
		//TODO CAche err
		return "", err
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
		//TODO CAche err
		return err
	}
	return nil
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

	// t, err := time.Parse("02/01/2006", req.ReleaseTime)
	// if err != nil {
	// 	return err
	// }
	rate, err := strconv.ParseFloat(req.Rate, 8)

	if err != nil {
		//TODO cache err
		return err
	}
	t := time.Now()
	movie := &Movie{
		ID:           id,
		Name:         req.Name,
		Rate:         rate,
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
func (s *Service) GetAllMovies(ctx context.Context, req FindRequest) ([]*types.MovieInfo, error) {

	if req.Offset > 50 {
		req.Offset = 50
	}
	movies, err := s.repository.GetAllMovies(ctx, req)
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

func (s *Service) GetMovieByName(name string) (*types.MovieInfo, error) {

	ctx := context.Background()
	movie, err := s.repository.GetMovieByName(ctx, name)
	if err != nil {
		//TODO CAche err
		return nil, err
	}
	return movie.ConvertMovieToMovieResponse(), nil
}

func (s *Service) GetMovieById(ctx context.Context, id string) (*types.MovieInfo, error) {
	movie, err := s.repository.GetMovieById(ctx, id)
	if err != nil {
		//TODO CAche err
		return nil, err
	}
	return movie.ConvertMovieToMovieResponse(), nil
}
