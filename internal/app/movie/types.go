package movie

import (
	"time"

	"github.com/dqkcode/movie-database/internal/app/types"
)

type (
	Movie struct {
		ID           string     `bson:"_id"`
		Name         string     `bson:"name"`
		Rate         float64    `bson:"rate"`
		Directors    []string   `bson:"directors"`
		Writers      []string   `bson:"writers"`
		TrailersPath []string   `bson:"trailers_path"`
		ImagesPath   []string   `bson:"images_path"`
		Casts        []string   `bson:"casts"`
		Genres       []string   `bson:"genres"`
		Storyline    string     `bson:"storyline"`
		UserReviews  []string   `bson:"user_reviews"`
		MovieLength  int        `bson:"movie_length"`
		ReleaseTime  string     `bson:"release_time"`
		CreatedAt    *time.Time `bson:"created_at"`
		UpdatedAt    *time.Time `bson:"updated_at"`
		UserId       string     `bson:"user_id"`
	}

	CreateRequest struct {
		Name        string   `validate:"required" json:"name"`
		Rate        float64  `validate:"required" json:"rate"`
		Directors   []string `validate:"required" json:"directors"`
		Writers     []string `validate:"required" json:"writers"`
		Casts       []string `validate:"required" json:"casts"`
		Genres      []string `validate:"required" json:"genres"`
		MovieLength int      `validate:"required" json:"movie_length"`
		ReleaseTime string   `validate:"required" json:"release_time"`
	}

	UpdateRequest struct {
		Name         string   `json:"name"`
		Rate         float64  `json:"rate"`
		Directors    []string `json:"directors"`
		Writers      []string `json:"writers"`
		TrailersPath []string `json:"trailers_path"`
		ImagesPath   []string `json:"images_path"`
		Genres       []string `json:"genres"`
		Casts        []string `json:"casts"`
		Storyline    string   `json:"storyline"`
		UserReviews  []string `json:"user_reviews"`
		MovieLength  int      `json:"movie_length"`
		ReleaseTime  string   `json:"release_time"`
	}
	FindRequest struct {
		Name        string   `json:"name"`
		Rate        float64  `json:"rate"`
		Directors   []string `json:"directors"`
		Writers     []string `json:"writers"`
		Genres      []string `json:"genres"`
		Casts       []string `json:"casts"`
		MovieLength int      `json:"movie_length"`
		ReleaseTime string   `json:"release_time"`
		CreatedByID string   `json:"created_by_id"`
		Offset      int      `json:"offset"`
		Limit       int      `json:"limit"`
		Selects     []string `json:"selects"`
		SortBy      []string `json:"sort_by"`
	}
)

func (m *Movie) ConvertMovieToMovieResponse() *types.MovieInfo {
	return &types.MovieInfo{
		ID:           m.ID,
		Name:         m.Name,
		Rate:         m.Rate,
		Directors:    m.Directors,
		Writers:      m.Writers,
		TrailersPath: m.TrailersPath,
		ImagesPath:   m.ImagesPath,
		Casts:        m.Casts,
		Genres:       m.Genres,
		Storyline:    m.Storyline,
		UserReviews:  m.UserReviews,
		MovieLength:  m.MovieLength,
		ReleaseTime:  m.ReleaseTime,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		UserId:       m.UserId,
	}
}
