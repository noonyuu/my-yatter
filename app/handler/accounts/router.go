package accounts

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
	ru usecase.Relationship
	au usecase.Account
	ar repository.Account
}

// Create Handler for `/v1/accounts/`
func NewRouter(ru usecase.Relationship, au usecase.Account, ar repository.Account) http.Handler {
	r := chi.NewRouter()

	h := &handler{
		ru: ru,
		au: au,
		ar: ar,
	}
	r.Post("/", h.Create)
	r.With(auth.Middleware(ar)).Post("/update_credentials", h.UpdateCredential)
	r.With(auth.Middleware(ar)).Get("/relationships", h.Relationships)
	r.Route("/{username}", func(r chi.Router) {
		r.Use(username)
		r.Get("/", h.FindByUsername)
		r.With(auth.Middleware(ar)).Post("/follow", h.Follow)
		r.Get("/following", h.Following)
		r.Get("/followers", h.Followers)

	})
	return r
}

func username(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		if username == "" {
			http.Error(w, "username is required", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), "username", username)
		n.ServeHTTP(w, r.WithContext(ctx))
	})
}
