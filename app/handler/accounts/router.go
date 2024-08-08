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
	rr repository.Relationship
	au usecase.Account
}

// Create Handler for `/v1/accounts/`
func NewRouter(rr repository.Relationship, au usecase.Account, ar repository.Account) http.Handler {
	r := chi.NewRouter()

	h := &handler{
		rr: rr,
		au: au,
	}
	r.Post("/", h.Create)
	r.Group(func(r chi.Router) {
		
		r.Use(auth.Middleware(ar))
		r.Post("/update_credentials", h.UpdateCredential)
	})
	r.Route("/{username}", func(r chi.Router) {
		r.Use(username)
		r.Get("/", h.FindByUsername)
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
