package user

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type (
	repository interface {
		Create(ctx context.Context, user User) (string, error)
		Update(ctx context.Context, user User) error
		FindUserByEmail(ctx context.Context, email string) (*User, error)
		Delete(ctx context.Context, id string) error
		CheckEmailIsRegisted(ctx context.Context, email string) error
	}
	Service struct {
		repository
	}
)

func NewService(repo repository) *Service {
	return &Service{
		repository: repo,
	}

}
func (s *Service) Register(ctx context.Context, req RegisterRequest) (string, error) {
	if err := validator.New().Struct(req); err != nil {
		logrus.Errorf("failed to validation, err: %v", err)
		return "", err
	}

	err := s.repository.CheckEmailIsRegisted(ctx, req.Email)
	if err == ErrDBQuery {
		return "", err
	} else if err == ErrUserAlreadyExist {
		return "", err
	}

	pwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("failed to generate password, err: %v", err)
		return "", ErrGenPasswordFailed
	}

	user := User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Gender:    req.Gender,
		Password:  string(pwd),
		CreatedAt: time.Now(),
	}
	id, err := s.repository.Create(ctx, user)
	if err != nil {
		logrus.Errorf("failed to create user, err: %v", err)
		return "", ErrCreateUserFailed
	}
	return id, nil
}

func (s *Service) Update(ctx context.Context, req UpdateInfoRequest) error {
	userUp := User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Gender:    req.Gender,
	}

	err := s.repository.Update(ctx, userUp)
	if err != nil && err == ErrUpdateUserFailed {
		logrus.Errorf("failed to update user, err: %v", err)
		return ErrUpdateUserFailed
	}
	return nil
}

func (s *Service) ChangePassword(ctx context.Context, req ChangePasswordRequest) error {
	return nil
}
func (s *Service) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.repository.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
