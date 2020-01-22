package auth

import (
	"net/http"

	"github.com/dqkcode/movie-database/internal/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {

	return []router.Route{
		{
			Path:    "/api/v1/login",
			Method:  http.MethodPost,
			Handler: h.Login,
		},
	}
}
