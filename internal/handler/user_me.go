package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/repo"
)

type UserMeHandler struct {
	repo     *repo.UserRepository
	pHandler *UserPreferencesHandler
}

func NewUserMeHandler(repo *repo.UserRepository, pHandler *UserPreferencesHandler) *UserMeHandler {
	return &UserMeHandler{
		repo:     repo,
		pHandler: pHandler,
	}
}

func (h *UserMeHandler) Routes(r chi.Router) {
	r.Route("/me", func(me chi.Router) {
		me.Put("/", func(w http.ResponseWriter, r *http.Request) {})
		me.Delete("/", func(w http.ResponseWriter, r *http.Request) {})

		me.Route("/email", func(email chi.Router) {
			email.Post("/request", func(w http.ResponseWriter, r *http.Request) {})
			email.Post("/confirm", func(w http.ResponseWriter, r *http.Request) {})
			email.Patch("/", func(w http.ResponseWriter, r *http.Request) {})
		})

		h.pHandler.Routes(me)
	})

}
