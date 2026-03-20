package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/repo"
)

type UserPreferencesHandler struct {
	repo *repo.UserPreferencesRepository
}

func NewUserPreferencesHandler(repo *repo.UserPreferencesRepository) *UserPreferencesHandler {
	return &UserPreferencesHandler{
		repo: repo,
	}
}

func (h *UserPreferencesHandler) Routes(r chi.Router) {
	// TODO: owner restrictions
	r.Route("/preferences", func(pref chi.Router) {
		pref.Post("/", func(w http.ResponseWriter, r *http.Request) {})
		pref.Put("/", func(w http.ResponseWriter, r *http.Request) {})
	})
}
