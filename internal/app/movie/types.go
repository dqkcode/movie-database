package movie

import (
	"time"

	"github.com/dqkcode/movie-database/internal/app/types"
)

type (
	Movie struct {
		ID           string    `bson:"_id"`
		Name         string    `bson:"name"`
		Rate         string    `bson:"rate"`
		Director     string    `bson:"director"`
		Writers      []string  `bson:"writers"`
		TrailersPath []string  `bson:"trailers_path"`
		ImagesPath   []string  `bson:"images_path"`
		Casts        []string  `bson:"casts"`
		Genres       []string  `bson:"genres"`
		Storyline    string    `bson:"storyline"`
		UserReviews  []string  `bson:"user_reviews"`
		MovieLength  int       `bson:"movie_length"`
		ReleaseTime  string    `bson:"release_time"`
		CreatedAt    time.Time `bson:"created_at"`
		UpdatedAt    time.Time `bson:"updated_at"`
		UserId       string    `bson:"user_id"`
	}

	CreateRequest struct {
		Name        string   `validate:"required" json:"name"`
		Director    string   `validate:"required" json:"director"`
		Writers     []string `validate:"required" json:"writers"`
		Casts       []string `validate:"required" json:"casts"`
		Genres      []string `validate:"required" json:"genres"`
		MovieLength int      `validate:"required" json:"movie_length"`
		ReleaseTime string   `validate:"required" json:"release_time"`
	}

	UpdateRequest struct {
		Name         string   `json:"name"`
		Rate         string   `json:"rate"`
		Director     string   `json:"director"`
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
)

func (m *Movie) ConvertMovieToMovieResponse() *types.MovieInfo {
	return &types.MovieInfo{
		ID:           m.ID,
		Name:         m.Name,
		Rate:         m.Rate,
		Director:     m.Director,
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
