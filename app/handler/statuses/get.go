package statuses

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"yatter-backend-go/app/usecase"
)

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	dto, err := h.su.FindByStatus(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newResponse(dto)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// レスポンス用の構造体
type HTTPResponse struct {
	ID      int
	Account struct {
		ID        int
		UserName  string
		CreatedAt time.Time
	}
	Content   string
	CreatedAt time.Time
}

type GetStatus struct {
	Statues *HTTPResponse
}

func newResponse(dto *usecase.GetStatusDTO) *GetStatus {
	return &GetStatus{
		Statues: &HTTPResponse{
			ID: dto.Status.ID,
			Account: struct {
				ID        int
				UserName  string
				CreatedAt time.Time
			}{
				ID:        int(dto.Account.ID),
				UserName:  dto.Account.Username,
				CreatedAt: dto.Account.CreateAt,
			},
			Content:   dto.Status.Content,
			CreatedAt: dto.Status.CreatedAt,
		},
	}
}
