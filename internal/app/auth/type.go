package auth

import (
	"github.com/dgrijalva/jwt-go"
)

type (
	LoginRequest struct {
		Email    string `validate:"required,email" json:"email"`
		Password string `validate:"required" json:"password"`
	}
	Claims struct {
		Email string `validate:"required,email" json:"email"`
		Id    string
		Role  string
		jwt.StandardClaims
	}
)
