package accounts

import (
	"encoding/json"
	"net/http"
	"strings"
	"yatter-backend-go/app/domain/auth"
	"yatter-backend-go/app/domain/object"
)

func (h *handler) Relationships(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account_info := auth.AccountOf(ctx) // 認証情報を取得する
	usernames := r.URL.Query().Get("username")
	if usernames == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}
	pairNames := strings.Split(usernames, ",")

	pairAccount := make([]*object.Account, 0, len(pairNames))
	for _, name := range pairNames {
		acc, err := h.ar.FindByUsername(ctx, name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pairAccount = append(pairAccount, acc)
	}

	rel, err := h.ru.GetRelationships(ctx, account_info, pairAccount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
