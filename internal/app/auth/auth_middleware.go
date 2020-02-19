package auth

import (
	"context"
	"errors"
	"net/http"

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
		claims, err := VerifyToken(tokenString)
		if errors.Is(err, ErrCompareToken) {
			logrus.Errorf("can not compare token, error: %v", err)
			types.ResponseJson(w, "", types.Normal().Internal)
			return
		}
		if errors.Is(err, ErrTokenInvalid) {
			logrus.Errorf("the token is invalid, error: %v", err)
			types.ResponseJson(w, "", types.Auth().TokenInvalid)
			return
		}
		NewUser := &types.UserInfo{
			ID:    claims.Id,
			Email: claims.Email,
			Role:  claims.Role,
		}
		newCtx := context.WithValue(r.Context(), types.UserContextKey, NewUser)
		r = r.WithContext(newCtx)
		h.ServeHTTP(w, r)
	})
}

func NewContextWithUser(ctx context.Context, u *types.UserInfo) context.Context {
	return context.WithValue(ctx, types.UserContextKey, u)
}

func GetRoleFromContext(ctx context.Context) types.Role {
	u := ctx.Value(types.UserContextKey).(*types.UserInfo)
	return (types.Role)(u.Role)
}
