package statuses

import (
	"context"
	"net/http"
	"yatter-backend-go/app/domain/repository"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/usecase"

	"github.com/go-chi/chi/v5"
)

// Implementation of handler
type handler struct {
	su usecase.Status
}

// Create Handler for `/v1/statuses/`
func NewRouter(ar repository.Account, su usecase.Status) http.Handler {
	r := chi.NewRouter()
	h := &handler{su: su}
	// r.Group()により、特定のグループに対してミドルウェアを適用する
	// グループに対して適用されたミドルウェアは、そのグループに属する全てのエンドポイントに対して適用される
	r.Group(func(r chi.Router) {
		// リクエストの認証を行う
		r.Use(auth.Middleware(ar))
		r.Post("/", h.Create)
		// TODO /{username}のハンドラを追加する
	})
	r.Route("/{id}", func(r chi.Router) {
		r.Use(id)
		r.Get("/", h.Get)
		r.With(auth.Middleware(ar)).Delete("/", h.Delete)
	})
	return r
}

func id(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), "id", id)
		n.ServeHTTP(w, r.WithContext(ctx))
	})
}
