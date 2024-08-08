package timeline

import (
	"net/http"
	"yatter-backend-go/app/domain/repository"
	"yatter-backend-go/app/usecase"

	"github.com/go-chi/chi/v5"
)

// Implementation of handler
type handler struct {
	tu usecase.Timeline
}

// ar repository.Account, mr repository.Media, sr repository.Status, tr repository.Timeline)
func NewRouter(tu usecase.Timeline, ar repository.Account) http.Handler {
	r := chi.NewRouter()
	h := &handler{tu: tu}

	r.Get("/public", h.GetPublicTimeline)

	return r
}
