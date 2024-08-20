package timeline

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/domain/auth"
)

func (h *handler) Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account_info := auth.AccountOf(ctx) // 認証情報を取得する
	if account_info == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	maxId := r.URL.Query().Get("max_id")
	sinceId := r.URL.Query().Get("since_id")
	limit := r.URL.Query().Get("limit")

	dto, err := h.tu.GetHomeTimeline(ctx, account_info.ID,maxId, sinceId, limit)
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
