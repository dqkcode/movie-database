package watchlist

import (
	"errors"
	"time"
)

type (
	Watchlist struct {
		ID        string    `bson:"_id"`
		UserID    string    `bson:"user_id"`
		Name      string    `bson:"name"`
		Share     bool      `bson:"share"`
		CreatedAt time.Time `bson:"created_at"`
		UpdatedAt time.Time `bson:"updated_at"`
	}
	WatchlistMovie struct {
		ID          string `bson:"_id"`
		WatchlistID string `bson:"watchlist_id"`
		MovieID     string `bson:"movie_id"`
	}
	WatchlistResponse struct {
		ID        string    `json:"_id"`
		UserID    string    `json:"user_id"`
		Name      string    `json:"name"`
		Share     bool      `json:"share"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

var (
	ErrWatchlistNotFound = errors.New("error watchlist not found")
	ErrPermissionDeny    = errors.New("error permission deny")
	ErrMovieNotFound     = errors.New("error movie not found")
	ErrDB                = errors.New("error DB")
)

func (w *Watchlist) ConvertToWatchlistResponse() *WatchlistResponse {
	return &WatchlistResponse{
		ID:        w.ID,
		Name:      w.Name,
		Share:     w.Share,
		UpdatedAt: w.UpdatedAt,
		CreatedAt: w.CreatedAt,
		UserID:    w.UserID,
	}
}
