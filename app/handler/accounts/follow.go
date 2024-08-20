package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yatter-backend-go/app/domain/auth"
)

func (h *handler) Follow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// ctx.Value()で、リクエストのコンテキストから値を取得する
	username := ctx.Value("username").(string)
	follower := auth.AccountOf(ctx) // 認証情報を取得する
	if follower == nil {
		panic(fmt.Sprintf("Must Implement Status Creation And Check Account Info %v", follower))
	}
	// usernameからaccountを取得
	followee, err := h.ar.FindByUsername(ctx, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.ru.FollowUser(ctx, follower, followee)
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
