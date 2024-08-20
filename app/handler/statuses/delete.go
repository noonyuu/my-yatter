package statuses

import (
	"net/http"
)

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// ctx.Value()で、リクエストのコンテキストから値を取得する
	id := ctx.Value("id").(string)

	err := h.su.DeleteStatus(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}
