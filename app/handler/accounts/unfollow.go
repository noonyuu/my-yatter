package accounts

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/domain/auth"
)

func (h *handler) UnFollow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// ctx.Value()で、リクエストのコンテキストから値を取得する
	username := ctx.Value("username").(string)
	account_info := auth.AccountOf(ctx) // 認証情報を取得する
	if account_info == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// usernameからaccountを取得
	followee, err := h.ar.FindByUsername(ctx, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.ru.UnFollowUser(ctx, account_info, followee)
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
