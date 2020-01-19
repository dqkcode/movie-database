package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dqkcode/movie-database/internal/app/types"
)

type (
	service interface {
		Login(ctx context.Context, req LoginRequest) (string, error)
	}
	Handler struct {
		srv service
	}
)

func NewHandler(srv service) *Handler {
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
	token, err := h.srv.Login(r.Context(), req)
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
