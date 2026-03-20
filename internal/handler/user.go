package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/repo"
)

type UserHandler struct {
	uRepo        *repo.UserRepository
	sRepo        *repo.SessionRepository
	uMeHandler   *UserMeHandler
	uPermHandler *UserPermissionHandler
	verification *VerificationHandler
}

func NewUserHandler(
	uRepo *repo.UserRepository,
	sRepo *repo.SessionRepository,
	uMeHandler *UserMeHandler,
	uPermHandler *UserPermissionHandler,
	verification *VerificationHandler,
) *UserHandler {
	return &UserHandler{
		uRepo:        uRepo,
		sRepo:        sRepo,
		verification: verification,
		uMeHandler:   uMeHandler,
		uPermHandler: uPermHandler,
	}
}

func (h *UserHandler) Routes(r chi.Router) {
	r.Route("/user", func(user chi.Router) {
		user.Get("/", func(w http.ResponseWriter, r *http.Request) {})
		user.Route("/:id", func(uid chi.Router) {
			uid.Get("/:id", func(w http.ResponseWriter, r *http.Request) {})
			uid.Group(func(admin chi.Router) {
				admin.Use(middlewares.SessionMiddleware(h.sRepo.GetByID))
				// root admin only routes, for now only check if the user stored in the session is admin
				// later this will need to implement permissions
				admin.Put("/:id", func(w http.ResponseWriter, r *http.Request) {})
				admin.Delete("/:id", func(w http.ResponseWriter, r *http.Request) {})
			})
		})

		user.Group(func(admin chi.Router) {
			admin.Use(middlewares.SessionMiddleware(h.sRepo.GetByID))
			// root admin only routes, for now only check if the user stored in the session is admin
			// later this will need to implement permissions
			admin.Post("/", func(w http.ResponseWriter, r *http.Request) {})
		})

		user.Group(func(auth chi.Router) {
			auth.Use(middlewares.SessionMiddleware(h.sRepo.GetByID))

			auth.Route("/verify", func(verify chi.Router) {
				verify.Post("/email", func(w http.ResponseWriter, r *http.Request) {})
			})

			h.uMeHandler.Routes(auth)
			h.uPermHandler.Routes(auth)
		})
	})
}
