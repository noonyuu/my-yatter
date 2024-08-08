package statuses

import (
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
	})
	r.Get("/{id}", h.Get)

	return r
}

func NewTimelineRouter(su usecase.Status) http.Handler {
	r := chi.NewRouter()
	h := &handler{su: su}
	
	r.Get("/public", h.GetPublicTimeline)

	return r
}