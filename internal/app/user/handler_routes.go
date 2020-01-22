package user

import (
	"net/http"

	"github.com/dqkcode/movie-database/internal/app/auth"

	"github.com/dqkcode/movie-database/internal/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {

	return []router.Route{
		{
			Handler: h.Register,
			Method:  http.MethodPost,
			Path:    "/api/v1/register",
		},
		{
			Handler:     h.Update,
			Method:      http.MethodPut,
			Path:        "/api/v1/update",
			Middlewares: []router.Middleware{auth.AuthMiddleware},
		},
	}
}
