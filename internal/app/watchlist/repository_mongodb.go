package watchlist

import (
	"errors"

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

func (r *MongoDBRepository) AddMovieToWatchlist(userID, movieID string) error {
	s := r.session.Clone()
	doc := bson.M{
		"user_id":  userID,
		"movie_id": movieID,
	}
	return r.getCollection(s).Insert(doc)
}

func (r *MongoDBRepository) GetWatchlistById(id string) (*Watchlist, error) {
	s := r.session.Clone()
	w := &Watchlist{}
	if err := r.getCollection(s).FindId(id).One(w); err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return nil, ErrWatchlistNotFound
		}
		return nil, err
	}
	return w, nil
}

func (r *MongoDBRepository) getCollection(s *mgo.Session) *mgo.Collection {
	return r.session.DB("").C("watchlist")
}
