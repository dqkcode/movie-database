package user

import (
	"context"
)

type (
	repository interface {
		// CURD
		Create(ctx context.Context, user User) (string, error)
		Update(ctx context.Context, user User) (string, error)
		Read(ctx context.Context, id string) (User, error)
		Delete(ctx context.Context, id string) error
	}
	Service struct {
		repository
	}
)

func NewService(repo repository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (string, error) {

	return "", nil
}

func (s *Service) RUpdate(ctx context.Context, req UpdateInfoRequest) (User, error) {
	return User{}, nil
}

func (s *Service) ChangePassword(ctx context.Context, req ChangePasswordRequest) error {
	return nil
}
