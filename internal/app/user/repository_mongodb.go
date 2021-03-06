package user

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
	c := m.getCollection(s)
	obj := bson.M{
		"updated_at": time.Now(),
	}

	if user.FirstName != "-" && user.FirstName != "" {
		obj["first_name"] = user.FirstName
	} else if user.FirstName == "-" {
		obj["first_name"] = ""
	}

	if user.LastName != "-" && user.LastName != "" {
		obj["last_name"] = user.LastName
	} else if user.LastName == "-" {
		obj["last_name"] = ""
	}

	if user.Gender >= 1 && user.Gender <= 3 {
		obj["gender"] = user.Gender
	} else {
		obj["gender"] = 3
	}

	err := c.UpdateId(user.ID, bson.M{
		"$set": obj,
	})
	if errors.Is(err, mgo.ErrNotFound) {
		return ErrUserNotFound
	}
	if err != nil {
		return err
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
	if errors.Is(err, mgo.ErrNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (m *MongoDBRepository) GetAllUsers(ctx context.Context) ([]*User, error) {
	s := m.session.Clone()
	defer s.Close()

	users := []*User{}
	err := m.getCollection(s).Find(nil).All(&users)

	if len(users) == 0 {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (m *MongoDBRepository) Delete(ctx context.Context, id string) error {
	s := m.session.Clone()
	defer s.Close()
	err := m.getCollection(s).RemoveId(id)
	if errors.Is(err, mgo.ErrNotFound) {
		return ErrUserNotFound
	}
	if err != nil {
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
	if errors.Is(err, mgo.ErrNotFound) {
		return ErrUserNotFound
	}
	if err != nil {
		return ErrDB
	}
	return ErrUserAlreadyExist
}
func (m *MongoDBRepository) getCollection(session *mgo.Session) *mgo.Collection {
	return session.DB("").C("users")
}
