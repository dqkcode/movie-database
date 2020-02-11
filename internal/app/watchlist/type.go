package watchlist

import (
	"errors"
)

type (
	Watchlist struct {
		ID       string   `bson:"_id"`
		UserID   string   `bson:"user_id"`
		MovieIDs []string `bson:"movie_ids"`
	}
)

var (
	ErrWatchlistNotFound = errors.New("error watchlist not found")
	ErrDB                = errors.New("error DB")
)
