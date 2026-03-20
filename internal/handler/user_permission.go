package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/repo"
)

type UserPermissionHandler struct {
	repo *repo.UserPermissionRepository
}

func NewUserPermissionHandler(repo *repo.UserPermissionRepository) *UserPermissionHandler {
	return &UserPermissionHandler{
		repo: repo,
	}
}

func (h *UserPermissionHandler) Routes(r chi.Router) {
	r.Route("/permissions", func(perms chi.Router) {
		perms.Get("/:id", func(w http.ResponseWriter, r *http.Request) {})

		// TODO: add system restrictions
		perms.Post("/", func(w http.ResponseWriter, r *http.Request) {})
		perms.Put("/", func(w http.ResponseWriter, r *http.Request) {})
		perms.Delete("/", func(w http.ResponseWriter, r *http.Request) {})

	})
}
