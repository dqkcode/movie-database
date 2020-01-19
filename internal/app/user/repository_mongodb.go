package user

import (
	"context"
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

func (m *MongoDBRepository) Update(ctx context.Context, user User) (string, error) {
	return "", nil
}

func (m *MongoDBRepository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	s := m.session.Clone()
	defer s.Close()
	selector := bson.M{
		"email": email,
	}
	user := &User{}
	if err := m.getCollection(s).Find(selector).One(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (m *MongoDBRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *MongoDBRepository) CheckEmailIsRegisted(ctx context.Context, email string) bool {
	s := m.session.Clone()
	defer s.Close()
	selector := bson.M{
		"email": email,
	}

	if err := m.getCollection(s).Find(selector).One(&User{}); err != nil {
		return false
	}

	return true
}
func (m *MongoDBRepository) getCollection(session *mgo.Session) *mgo.Collection {
	return session.DB("").C("users")
}
