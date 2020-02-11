package api

import (
	"github.com/dqkcode/movie-database/internal/app/watchlist"
	"github.com/globalsign/mgo"
)

func NewWatchlistHandler(session *mgo.Session) *watchlist.Handler {
	repo := watchlist.NewRepository(session)
	watchlistSvc := watchlist.NewService(repo)
	return watchlist.NewHandler(watchlistSvc)
}
