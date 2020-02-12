package watchlist

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type MongoDBRepository struct {
	session *mgo.Session
}

func NewRepository(session *mgo.Session) *MongoDBRepository {
	return &MongoDBRepository{
		session: session,
	}
}

func (r *MongoDBRepository) CreateWatchlist(ctx context.Context, usrID, name string) (string, error) {
	s := r.session.Clone()
	defer s.Close()
	id := uuid.New().String()
	t := time.Now()
	doc := bson.M{
		"_id":        id,
		"user_id":    usrID,
		"name":       name,
		"share":      false,
		"created_at": t,
		"updated_at": t,
	}
	if err := r.getCollection("watchlist", s).Insert(doc); err != nil {
		return "", err
	}
	return id, nil
}

func (r *MongoDBRepository) AddMovieToWatchlist(movieID, watchlistID string) error {
	s := r.session.Clone()
	defer s.Close()

	doc := bson.M{
		"watchlist_id": watchlistID,
		"movie_id":     movieID,
	}
	if err := r.getCollection("watchlist_movie", s).Insert(doc); err != nil {
		return err
	}
	return nil
}
func (r *MongoDBRepository) DeleteMovieInWatchlist(movieID, watchlistID string) error {
	s := r.session.Clone()
	defer s.Close()

	doc := bson.M{
		"watchlist_id": watchlistID,
		"movie_id":     movieID,
	}
	if err := r.getCollection("watchlist_movie", s).Remove(doc); err != nil {
		return err
	}
	return nil
}

func (r *MongoDBRepository) GetWatchlistById(id string) (*Watchlist, error) {
	s := r.session.Clone()
	defer s.Close()
	w := &Watchlist{}
	if err := r.getCollection("watchlist", s).FindId(id).One(w); err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return nil, ErrWatchlistNotFound
		}
		return nil, err
	}
	return w, nil
}

func (r *MongoDBRepository) ListAllMovies(id string) ([]*WatchlistMovie, error) {
	s := r.session.Clone()
	defer s.Close()
	list := []*WatchlistMovie{}
	if err := r.getCollection("watchlist_movie", s).Find(bson.M{"watchlist_id": id}).All(&list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *MongoDBRepository) getCollection(name string, s *mgo.Session) *mgo.Collection {
	return r.session.DB("").C(name)
}
