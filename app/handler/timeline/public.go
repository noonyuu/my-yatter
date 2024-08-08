package timeline

import (
	"encoding/json"
	"net/http"
	"time"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/usecase"
)

func (h *handler) Public(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// 各パラメータを取得
	// onlyMedia := r.URL.Query().Get("only_media") // 未使用
	sinceId := r.URL.Query().Get("since_id")
	limit := r.URL.Query().Get("limit")

	// sta := object.NewStatus(status)

	dto, err := h.tu.GetPublicTimeline(ctx, sinceId, limit)
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
	ID        int
	Account   object.Account
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
			ID:        dto.Status[i].ID,
			Account:   *dto.Account[i],
			Content:   dto.Status[i].Content,
			CreatedAt: dto.
			Status[i].CreatedAt,
		}
	}

	return &GetTimelineStatus{
		Statuses: statuses,
	}
}
