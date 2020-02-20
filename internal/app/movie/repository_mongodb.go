package movie

import (
	"context"
	"errors"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type (
	MongoDBRepository struct {
		session *mgo.Session
	}
)

func NewMongoDBRepository(s *mgo.Session) *MongoDBRepository {
	return &MongoDBRepository{
		session: s,
	}
}

func (m *MongoDBRepository) Create(ctx context.Context, movie Movie) (string, error) {
	s := m.session.Clone()
	defer s.Close()
	if err := m.getCollection(s).Insert(movie); err != nil {
		return "", err
	}
	return movie.ID, nil
}

func (m *MongoDBRepository) DeleteById(ctx context.Context, id string) error {
	s := m.session.Clone()
	defer s.Close()

	err := m.getCollection(s).RemoveId(id)
	if errors.Is(err, mgo.ErrNotFound) {
		return ErrMovieNotFound
	}
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDBRepository) GetAllMovies(ctx context.Context, req FindRequest) ([]*Movie, error) {
	s := m.session.Clone()
	defer s.Close()
	r := bson.M{}
	if len(req.Genres) > 0 {
		r["genres"] = bson.M{
			"$in": req.Genres,
		}
	}
	if len(req.Casts) > 0 {
		r["casts"] = bson.M{
			"$in": req.Casts,
		}
	}
	if len(req.Writers) > 0 {
		r["writers"] = bson.M{
			"$in": req.Writers,
		}
	}
	if len(req.Directors) > 0 {
		r["directors"] = bson.M{
			"$in": req.Directors,
		}
	}

	if req.Name != "" {
		r["name"] = bson.RegEx{
			Pattern: req.Name,
			Options: "i",
		}
	}
	if req.Rate > 0 {
		r["rate"] = bson.M{
			"$gt": req.Rate,
		}
	}
	if req.CreatedByID != "" {
		r["user_id"] = bson.M{
			"$in": req.CreatedByID,
		}
	}
	if req.MovieLength > 0 {
		r["movie_length"] = bson.M{
			"$lt": req.MovieLength,
		}
	}
	if req.ReleaseTime != "" {
		r["release_time"] = bson.M{
			"$in": req.ReleaseTime,
		}
	}

	selects := bson.M{}
	if len(req.Selects) > 0 {
		for _, s := range req.Selects {
			selects[s] = 1
		}
	}
	if len(req.SortBy) == 1 && req.SortBy[0] == "" {
		req.SortBy[0] = "name"
	}

	movies := []*Movie{}
	err := m.getCollection(s).Find(r).Select(selects).Skip(req.Offset).Limit(req.Limit).Sort(req.SortBy...).All(&movies)
	if err != nil {
		return nil, err
	}
	if len(movies) == 0 {
		return nil, ErrMovieNotFound
	}
	return movies, nil
}

func (m *MongoDBRepository) GetMovieByName(ctx context.Context, name string) (*Movie, error) {
	s := m.session.Clone()
	defer s.Close()
	movie := &Movie{}
	err := m.getCollection(s).Find(bson.M{"name": name}).One(movie)
	if errors.Is(err, mgo.ErrNotFound) {
		return nil, ErrMovieNotFound
	}
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (m *MongoDBRepository) GetMovieById(ctx context.Context, id string) (*Movie, error) {
	s := m.session.Clone()
	defer s.Close()

	movie := &Movie{}
	err := m.getCollection(s).FindId(id).One(movie)
	if errors.Is(err, mgo.ErrNotFound) {
		return nil, ErrMovieNotFound
	}
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (m *MongoDBRepository) Update(ctx context.Context, movie *Movie) error {
	s := m.session.Clone()
	defer s.Close()

	obj := bson.M{
		"updated_at": time.Now(),
	}
	if movie.Name != "" {
		obj["name"] = movie.Name
	}
	if movie.Storyline != "" {
		obj["storyline"] = movie.Storyline
	}
	if movie.ReleaseTime != "" {
		obj["release_time"] = movie.ReleaseTime
	}
	if movie.Rate > -1 {
		obj["rate"] = movie.Rate
	}
	if movie.MovieLength > 0 {
		obj["movie_length"] = movie.MovieLength
	}
	if len(movie.Directors) > 0 {
		obj["directors"] = movie.Directors
	}
	if len(movie.Writers) > 0 {
		obj["writers"] = movie.Writers
	}
	if len(movie.Casts) > 0 {
		obj["casts"] = movie.Casts
	}
	if len(movie.UserReviews) > 0 {
		obj["user_reviews"] = movie.UserReviews
	}
	if len(movie.TrailersPath) > 0 {
		obj["trailers_path"] = movie.TrailersPath
	}
	if len(movie.ImagesPath) > 0 {
		obj["images_path"] = movie.ImagesPath
	}

	err := m.getCollection(s).UpdateId(movie.ID, bson.M{
		"$set": obj})
	if errors.Is(err, mgo.ErrNotFound) {
		return ErrMovieNotFound
	}
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDBRepository) getCollection(session *mgo.Session) *mgo.Collection {
	return session.DB("").C("movies")
}
