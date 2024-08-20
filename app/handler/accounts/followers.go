package accounts

import (
	"encoding/json"
	"net/http"
)

func (h *handler) Followers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// ctx.Value()で、リクエストのコンテキストから値を取得する
	username := ctx.Value("username").(string)
	sinceId := r.URL.Query().Get("since_id")
	limit := r.URL.Query().Get("limit")
	// usernameからaccountを取得
	account, err := h.ar.FindByUsername(ctx, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	follower, err := h.au.FollowerAccount(ctx, account, limit,sinceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(follower); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}