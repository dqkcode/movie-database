package types

import (
	"errors"
	"time"
)

type (
	MovieInfo struct {
		ID           string     `json:"_id,omitempty"`
		Name         string     `json:"name,omitempty"`
		Rate         float64    `json:"rate,omitempty"`
		Directors    []string   `json:"directors,omitempty"`
		Writers      []string   `json:"writers,omitempty"`
		TrailersPath []string   `json:"trailers_path,omitempty"`
		ImagesPath   []string   `json:"images_path,omitempty"`
		Casts        []string   `json:"casts,omitempty"`
		Genres       []string   `json:"genres,omitempty"`
		Storyline    string     `json:"storyline,omitempty"`
		UserReviews  []string   `json:"user_reviews,omitempty"`
		MovieLength  int        `json:"movie_length,omitempty"`
		ReleaseTime  string     `json:"release_time,omitempty"`
		CreatedAt    *time.Time `json:"created_at,omitempty"`
		UpdatedAt    *time.Time `json:"updated_at,omitempty"`
		UserId       string     `json:"user_id,omitempty"`
	}
)

var ErrMovieNotFound = errors.New("Movie not found")
