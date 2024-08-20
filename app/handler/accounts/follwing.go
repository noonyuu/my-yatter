package accounts

import (
	"encoding/json"
	"net/http"
)

func (h *handler) Following(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// ctx.Value()で、リクエストのコンテキストから値を取得する
	username := ctx.Value("username").(string)
	limit := r.URL.Query().Get("limit")
	// usernameからaccountを取得
	account, err := h.ar.FindByUsername(ctx, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	followee, err := h.au.FolloweeAccount(ctx, account, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(followee); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}