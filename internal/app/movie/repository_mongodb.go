package movie

import (
	"context"
	"time"

	"github.com/dqkcode/movie-database/internal/app/types"
	"github.com/globalsign/mgo/bson"

	"github.com/globalsign/mgo"
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

	if err := m.getCollection(s).RemoveId(id); err != nil {
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
	//TODO check
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
		r["user_id"] = req.CreatedByID
	}
	if req.MovieLength > 0 {
		r["movie_length"] = bson.M{
			"$lt": req.MovieLength,
		}
	}
	if req.ReleaseTime != "" {
		r["release_time"] = req.ReleaseTime
	}

	selects := bson.M{}
	if len(req.Selects) > 0 {
		for _, s := range req.Selects {
			selects[s] = 1
		}
	}
	if len(req.SortBy) == 1 {
		if req.SortBy[0] == "" {
			req.SortBy[0] = "name"
		}
	}
	movies := []*Movie{}
	err := m.getCollection(s).Find(r).Select(selects).Skip(req.Offset).Limit(req.Limit).Sort(req.SortBy...).All(&movies)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrMovieNotFound
		}
		return nil, err
	}
	if len(movies) == 0 {
		return nil, ErrMovieNotFound
	}
	return movies, nil
}

func (m *MongoDBRepository) GetAllMoviesByUserId(ctx context.Context, id string) ([]*Movie, error) {
	s := m.session.Clone()
	defer s.Close()

	movies := []*Movie{}
	err := m.getCollection(s).Find(bson.M{"user_id": id}).All(&movies)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrMovieNotFound
		}
		return nil, err
	}
	return movies, nil
}

func (m *MongoDBRepository) GetMovieByName(ctx context.Context, name string) (*Movie, error) {
	s := m.session.Clone()
	defer s.Close()

	movie := &Movie{}
	err := m.getCollection(s).Find(bson.M{"name": name}).One(movie)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, types.ErrMovieNotFound
		}
		return nil, err
	}
	return movie, nil
}

func (m *MongoDBRepository) GetMovieById(ctx context.Context, id string) (*Movie, error) {
	s := m.session.Clone()
	defer s.Close()

	movie := &Movie{}
	err := m.getCollection(s).FindId(id).One(movie)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrMovieNotFound
		}
		return nil, err
	}
	return movie, nil
}

func (m *MongoDBRepository) Update(ctx context.Context, movie *Movie) error {
	s := m.session.Clone()
	defer s.Close()

	err := m.getCollection(s).UpdateId(movie.ID, bson.M{
		"$set": bson.M{
			"name":          movie.Name,
			"rate":          movie.Rate,
			"directors":     movie.Directors,
			"writers":       movie.Writers,
			"trailers_path": movie.TrailersPath,
			"images_path":   movie.ImagesPath,
			"casts":         movie.Casts,
			"storyline":     movie.Storyline,
			"user_reviews":  movie.UserReviews,
			"movie_length":  movie.MovieLength,
			"release_time":  movie.ReleaseTime,
			"updated_at":    time.Now(),
		}})
	if err != nil {
		if err == mgo.ErrNotFound {
			return ErrMovieNotFound
		}

		return ErrUpdateMovieFailed
	}
	return nil
}

func (m *MongoDBRepository) getCollection(session *mgo.Session) *mgo.Collection {
	return session.DB("").C("movies")
}
