package types

import (
	"time"
)

type (
	MovieInfo struct {
		ID           string    `json:"_id"`
		Name         string    `json:"name"`
		Rate         string    `json:"rate"`
		Director     string    `json:"director"`
		Writers      []string  `json:"writers"`
		TrailersPath []string  `json:"trailers_path"`
		ImagesPath   []string  `json:"images_path"`
		Casts        []string  `json:"casts"`
		Genres       []string  `json:"genres"`
		Storyline    string    `json:"storyline"`
		UserReviews  []string  `json:"user_reviews"`
		MovieLength  int       `json:"movie_length"`
		ReleaseTime  string    `json:"release_time"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		UserId       string    `json:"user_id"`
	}
)
