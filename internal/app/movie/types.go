package movie

import (
	"time"

	"github.com/dqkcode/movie-database/internal/app/types"
)

type (
	Movie struct {
		ID           string    `bson:"_id"`
		Name         string    `bson:"name"`
		Rate         int8      `bson:"rate"`
		Description  string    `bson:"description"`
		Director     string    `bson:"director"`
		Writers      []string  `bson:"writers"`
		Stars        []string  `bson:"stars"`
		TrailersPath []string  `bson:"trailers_path"`
		ImagesPath   []string  `bson:"images_path"`
		Casts        []string  `bson:"casts"`
		Genres       []string  `bson:"genres"`
		Storyline    string    `bson:"storyline"`
		UserReviews  []string  `bson:"user_reviews"`
		MovieLength  int       `bson:"movie_length"`
		ReleaseTime  time.Time `bson:"release_time"`
		CreatedAt    time.Time `bson:"created_at"`
		UpdatedAt    time.Time `bson:"updated_at"`
		UserId       string    `bson:"user_id"`
	}

	CreateRequest struct {
		Name        string   `validate:"required" json:"name"`
		Description string   `validate:"required" json:"description"`
		Director    string   `validate:"required" json:"director"`
		Writers     []string `validate:"required" json:"writers"`
		Stars       []string `validate:"required" json:"stars"`
		Casts       []string `validate:"required" json:"casts"`
		Genres      []string `validate:"required" json:"genres"`
		MovieLength int      `validate:"required" json:"movie_length"`
		ReleaseTime string   `validate:"required" json:"release_time"`
	}

	UpdateRequest struct {
		Name         string   `json:"name"`
		Rate         int8     `json:"rate"`
		Description  string   `json:"description"`
		Director     string   `json:"director"`
		Writers      []string `json:"writers"`
		Stars        []string `json:"stars"`
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
		Description:  m.Description,
		Director:     m.Director,
		Writers:      m.Writers,
		Stars:        m.Stars,
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
