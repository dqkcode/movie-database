package user

import (
	"net/http"

	"github.com/dqkcode/movie-database/internal/app/auth"

	"github.com/dqkcode/movie-database/internal/pkg/http/router"
)

const version1 = "/api/v1"

func (h *Handler) Routes() []router.Route {

	return []router.Route{
		{
			Handler: h.Register,
			Method:  http.MethodPost,
			Path:    version1 + "/register",
		},
		{
			Handler:     h.Update,
			Method:      http.MethodPut,
			Path:        version1 + "/users",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
		},
		{
			Handler:     h.DeleteUser,
			Method:      http.MethodDelete,
			Path:        version1 + "/users/{id}",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
		},
		{
			Handler:     h.FindUserById,
			Method:      http.MethodGet,
			Path:        version1 + "/users/{id}",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
		},
		{
			Handler:     h.GetAllUsers,
			Method:      http.MethodGet,
			Path:        version1 + "/users",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
		},
	}
}
