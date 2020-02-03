package auth

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dqkcode/movie-database/internal/app/types"
	"github.com/sirupsen/logrus"
)

func AuthMiddleware(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			types.ResponseJson(w, "", types.Auth().Unauthorized)
			return
		}
		var tokenValidType string
		if tokenString[:7] == "Bearer " {
			tokenValidType = tokenString[7:]
		} else {
			tokenValidType = tokenString
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenValidType, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("my_secret_key"), nil
		})
		if err != nil {
			logrus.Errorf("Can not compare token, error: %v", err)
			types.ResponseJson(w, "", types.Auth().TokenInvalid)
			return

		}
		if !token.Valid {
			logrus.Errorf("the token is invalid, error: %v", err)
			return
		}

		NewUser := &types.UserInfo{
			ID:    claims.Id,
			Email: claims.Email,
			Role:  claims.Role,
			//TODO add some info
		}
		newCtx := context.WithValue(r.Context(), "user", NewUser)
		r = r.WithContext(newCtx)
		h.ServeHTTP(w, r)
	})

}
