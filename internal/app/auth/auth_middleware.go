package auth

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dqkcode/movie-database/internal/app/types"
	"github.com/dqkcode/movie-database/internal/app/user"
	"github.com/sirupsen/logrus"
)

// func GetUserInfoMiddleware() func(http.Handler) http.Handler {
// 	return func(inner http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			tokenString := r.Header.Get("Authorization")
// 			if tokenString == "" {
// 				json.NewEncoder(w).Encode(types.Response{
// 					Code:  types.CodeFail,
// 					Error: "Unauthorized",
// 				})
// 				return
// 			}
// 			tokenValidType := strings.Replace(tokenString, "Bearer ", "", 7)
// 			claims := &Claims{}
// 			token, err := jwt.ParseWithClaims(tokenValidType, claims, func(token *jwt.Token) (interface{}, error) {
// 				return []byte("my_secret_key"), nil
// 			})
// 			if err != nil {
// 				logrus.Errorf("Can not compare token, error: %v", err)
// 				return

// 			}
// 			if !token.Valid {
// 				logrus.Errorf("the token is invalid, error: %v", err)
// 				return
// 			}

// 			NewUser := &user.User{
// 				ID:    claims.Id,
// 				Email: claims.Email,
// 			}
// 			newCtx := context.WithValue(r.Context(), "user", NewUser)
// 			r = r.WithContext(newCtx)
// 			inner.ServeHTTP(w, r)
// 		})
// 	}

// }

func UserInfoMiddleware(h http.Handler) http.Handler {

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
			return

		}
		if !token.Valid {
			logrus.Errorf("the token is invalid, error: %v", err)
			return
		}

		NewUser := &user.User{
			ID:    claims.Id,
			Email: claims.Email,
		}
		newCtx := context.WithValue(r.Context(), "user", NewUser)
		r = r.WithContext(newCtx)
		h.ServeHTTP(w, r)
	})

}
