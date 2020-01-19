package auth

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dqkcode/movie-database/internal/app/user"
	"github.com/go-playground/validator/v10"

	"github.com/sirupsen/logrus"
)

type (
	Service struct {
	}
)

func NewService() *Service {
	return &Service{}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (string, error) {

	if err := validator.New().Struct(req); err != nil {
		logrus.Errorf("failed to validation, err: %v", err)
		return "", err
	}

	user, err := user.Service.
	if err != nil {
		return "", err
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Email: req.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString("jwtKey")
	if err != nil {
		logrus.Errorf("Signing string fail")
		return "", err
	}

	return tokenString, nil
}
