package types

import "time"

type (
	MovieInfo struct {
		ID           string    `json:"_id"`
		Name         string    `json:"name"`
		Rate         int8      `json:"rate"`
		Description  string    `json:"description"`
		Director     string    `json:"director"`
		Writers      []string  `json:"writers"`
		Stars        []string  `json:"stars"`
		TrailersPath []string  `json:"trailers_path"`
		ImagesPath   []string  `json:"images_path"`
		Casts        []string  `json:"casts"`
		Genres       []string  `json:"genres"`
		Storyline    string    `json:"storyline"`
		UserReviews  []string  `json:"user_reviews"`
		MovieLength  int       `json:"movie_length"`
		ReleaseTime  time.Time `json:"release_time"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		UserId       string    `json:"user_id"`
	}
)
