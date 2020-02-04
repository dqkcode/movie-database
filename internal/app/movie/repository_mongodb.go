package movie

import (
	"context"
	"time"

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

func (m *MongoDBRepository) GetAllMovies(ctx context.Context) ([]*Movie, error) {
	s := m.session.Clone()
	defer s.Close()

	movies := []*Movie{}
	err := m.getCollection(s).Find(nil).All(&movies)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrMovieNotFound
		}
		return nil, err
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
			"description":   movie.Description,
			"director":      movie.Director,
			"writers":       movie.Writers,
			"stars":         movie.Stars,
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
		return ErrUpdateMovieFailed
	}
	return nil
}

func (m *MongoDBRepository) getCollection(session *mgo.Session) *mgo.Collection {
	return session.DB("").C("movies")
}
