package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dqkcode/movie-database/internal/app/types"
)

type (
	service interface {
		Auth(ctx context.Context, email string, password string) (string, error)
	}
	Handler struct {
		srv service
	}
)

func NewHanler(srv service) *Handler {
	return &Handler{
		srv,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := h.srv.Auth(r.Context(), req.Email, req.Password)
	if err != nil {
		json.NewEncoder(w).Encode(types.Response{
			Code:  types.CodeFail,
			Error: err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(types.Response{
		Code: types.CodeSuccess,
		Data: map[string]interface{}{
			"token": token,
		},
	})
}
