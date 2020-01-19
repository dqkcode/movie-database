package auth

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dqkcode/movie-database/internal/app/api"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

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

	userSrv := api.NewUserServive()
	user, err := userSrv.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logrus.Errorf("password wrong")
		return "", err
	}
	expirationTime := time.Now().Add(500 * time.Minute)
	claims := &Claims{
		Email: req.Email,
		Id:    user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString([]byte("my_secret_key"))
	if err != nil {
		logrus.Errorf("Signing string fail")
		return "", err
	}

	return tokenString, nil
}
