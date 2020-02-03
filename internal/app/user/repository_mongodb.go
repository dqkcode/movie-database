package user

import (
	"context"
	"time"

	"github.com/dqkcode/movie-database/internal/app/types"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type (
	MongoDBRepository struct {
		session *mgo.Session
	}
)

func NewMongoDBRepository(session *mgo.Session) *MongoDBRepository {
	return &MongoDBRepository{
		session,
	}
}

func (m *MongoDBRepository) Create(ctx context.Context, user User) (string, error) {
	s := m.session.Clone()
	defer s.Close()
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	if err := m.getCollection(s).Insert(user); err != nil {
		return "", err
	}
	return user.ID, nil
}

func (m *MongoDBRepository) Update(ctx context.Context, user User) error {
	s := m.session.Clone()
	defer s.Close()
	u := ctx.Value("user").(*types.UserInfo)

	c := m.getCollection(s)
	err := c.UpdateId(u.ID, bson.M{
		"$set": bson.M{
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"gender":     user.Gender,
			"updated_at": time.Now(),
		}})
	if err != nil {
		return ErrUpdateUserFailed
	}

	return nil
}

func (m *MongoDBRepository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	s := m.session.Clone()
	defer s.Close()
	selector := bson.M{
		"email": email,
	}
	user := &User{}
	if err := m.getCollection(s).Find(selector).One(user); err != nil {
		return nil, err
	}
	return user, nil
}
func (m *MongoDBRepository) FindUserById(ctx context.Context, id string) (*User, error) {
	s := m.session.Clone()
	defer s.Close()

	user := &User{}
	err := m.getCollection(s).FindId(id).One(user)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
func (m *MongoDBRepository) GetAllUsers(ctx context.Context) ([]*User, error) {
	s := m.session.Clone()
	defer s.Close()

	users := []*User{}
	err := m.getCollection(s).Find(nil).All(&users)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return users, nil
}

func (m *MongoDBRepository) Delete(ctx context.Context, id string) error {
	s := m.session.Clone()
	defer s.Close()
	if err := m.getCollection(s).RemoveId(id); err != nil {
		return err
	}
	return nil
}

func (m *MongoDBRepository) CheckEmailIsRegisted(ctx context.Context, email string) error {
	s := m.session.Clone()
	defer s.Close()
	selector := bson.M{
		"email": email,
	}

	err := m.getCollection(s).Find(selector).One(&User{})
	if err != nil {
		if err == mgo.ErrNotFound {
			return ErrUserNotFound

		}
		return ErrDBQuery
	}
	return ErrUserAlreadyExist
}
func (m *MongoDBRepository) getCollection(session *mgo.Session) *mgo.Collection {
	return session.DB("").C("users")
}
