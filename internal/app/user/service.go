package user

import (
	"context"
	"errors"
	"time"

	"github.com/dqkcode/movie-database/internal/app/types"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type (
	policyService interface {
		Validate(ctx context.Context, obj, act string) bool
	}
	repository interface {
		Create(ctx context.Context, user User) (string, error)
		Update(ctx context.Context, user User) error
		FindUserByEmail(ctx context.Context, email string) (*User, error)
		FindUserById(ctx context.Context, id string) (*User, error)
		Delete(ctx context.Context, id string) error
		CheckEmailIsRegisted(ctx context.Context, email string) error
		GetAllUsers(ctx context.Context) ([]*User, error)
	}
	Service struct {
		repository repository
		policy     policyService
	}
)

func NewService(repo repository, policy policyService) *Service {
	return &Service{
		repository: repo,
		policy:     policy,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (string, error) {
	if err := validator.New().Struct(req); err != nil {
		logrus.Errorf("failed to validation, err: %v", err)
		return "", err
	}
	err := s.repository.CheckEmailIsRegisted(ctx, req.Email)

	if errors.Is(err, ErrDB) {
		return "", err
	}
	if errors.Is(err, ErrUserAlreadyExist) {
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
		Role:      "normal",
	}
	id, err := s.repository.Create(ctx, user)
	if err != nil {
		logrus.Errorf("failed to create user, err: %v", err)
		return "", ErrCreateUserFailed
	}
	return id, nil
}

func (s *Service) Update(ctx context.Context, id string, req UpdateInfoRequest) error {
	userCtx := ctx.Value(types.UserContextKey).(*types.UserInfo)
	if userCtx.ID != id && userCtx.Role != "admin" {
		return ErrPermissionDeny
	}
	userUp := User{
		ID:        id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Gender:    req.Gender,
	}
	err := s.repository.Update(ctx, userUp)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ChangePassword(ctx context.Context, req ChangePasswordRequest) error {
	return nil
}
func (s *Service) FindUserByEmail(ctx context.Context, email string) (*types.UserInfo, error) {

	user, err := s.repository.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user.ConvertUserToUserInfo(), nil
}

func (s *Service) FindUserById(ctx context.Context, id string) (*types.UserInfo, error) {
	if s.policy.Validate(ctx, types.PolicyObjectUser, types.PolicyActionRead) == false {
		return nil, ErrPermissionDeny
	}
	user, err := s.repository.FindUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user.TrimSecrectUserInfo().ConvertUserToUserInfo(), nil
}

func (s *Service) DeleteUser(ctx context.Context, id string) error {
	if s.policy.Validate(ctx, types.PolicyObjectUser, types.PolicyActionRead) == false {
		return ErrPermissionDeny
	}
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) GetAllUsers(ctx context.Context) ([]*types.UserInfo, error) {
	if s.policy.Validate(ctx, types.PolicyObjectUser, types.PolicyActionRead) == false {
		return nil, ErrPermissionDeny
	}
	users, err := s.repository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	var us []*types.UserInfo
	for _, u := range users {
		us = append(us, u.TrimSecrectUserInfo().ConvertUserToUserInfo())
	}
	return us, nil
}
