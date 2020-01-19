package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dqkcode/movie-database/internal/app/types"
)

type (
	service interface {
		Register(ctx context.Context, req RegisterRequest) (string, error)
		Update(ctx context.Context, req UpdateInfoRequest) error
		// ChangePassword(ctx context.Context, req ChangePasswordRequest) error
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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.srv.Register(r.Context(), req)
	if err != nil {
		json.NewEncoder(w).Encode(types.Response{
			Code:  types.CodeFail,
			Data:  "",
			Error: err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(types.Response{
		Code: types.CodeSuccess,
		Data: map[string]interface{}{
			"id": id,
		},
		Error: "",
	})
}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req UpdateInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.srv.Update(r.Context(), req)
	if err != nil {
		json.NewEncoder(w).Encode(types.Response{
			Code:  types.CodeFail,
			Data:  "",
			Error: err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(types.Response{
		Code:  types.CodeSuccess,
		Data:  r.Context().Value("user").(*User).ID,
		Error: "",
	})
}
