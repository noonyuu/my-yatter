package statuses

import (
	"encoding/json"
	"fmt"
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

	dto, err := h.statusUsecase.FindByStatus(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) GetPublicTimeline(w http.ResponseWriter, r *http.Request) {
	fmt.Print("GetPublicTimeline")
	ctx := r.Context()
	// 各パラメータを取得
	// onlyMedia := r.URL.Query().Get("only_media") // 未使用
	maxId := r.URL.Query().Get("max_id")
	sinceId := r.URL.Query().Get("since_id")
	limit := r.URL.Query().Get("limit")

	dto, err := h.statusUsecase.GetPublicTimeline(ctx, maxId, sinceId, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseStatus(dto)); err != nil {
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

type GetTimelineStatus struct {
	Statuses []*HTTPResponse
}

func responseStatus(dto *usecase.GetPublicStatusDTO) *GetTimelineStatus {
	statuses := make([]*HTTPResponse, len(dto.Status))
	for i := range dto.Status {
		statuses[i] = &HTTPResponse{
			ID: dto.Status[i].ID,
			Account: struct {
				ID        int
				UserName  string 
				CreatedAt time.Time
			}{
				ID:        int(dto.Account[i].ID),
				UserName:  dto.Account[i].Username,
				CreatedAt: dto.Account[i].CreateAt,
			},
			Content:   dto.Status[i].Content,
			CreatedAt: dto.Status[i].CreatedAt,
		}
	}

	return &GetTimelineStatus{
		Statuses: statuses,
	}
}
