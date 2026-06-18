package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/api/core/permission"
	"github.com/juevigrace/diva-server/internal/api/core/session"
	"github.com/juevigrace/diva-server/internal/api/core/user"
	"github.com/juevigrace/diva-server/internal/api/core/verification"
	"github.com/juevigrace/diva-server/internal/api/middlewares"
)

type AuthModule struct {
	Handler  *AuthHandler
	Repo  *AuthRepo
	sRepo *session.SessionRepo
	uRepo *user.UserRepo
}

func NewAuthModule(
	pRepo *permission.PermissionRepo,
	sRepo *session.SessionRepo,
	uRepo *user.UserRepo,
	vRepo *verification.VerificationRepo,
) *AuthModule {
	repo := NewAuthRepo(pRepo, sRepo, uRepo, vRepo)
	return &AuthModule{
		Handler:  NewAuthHandler(repo),
		Repo:  repo,
		sRepo: sRepo,
		uRepo: uRepo,
	}
}

func (m *AuthModule) Routes(r chi.Router) {
	r.Route("/auth", func(auth chi.Router) {
		auth.Post("/signIn", m.Handler.signIn)
		auth.Post("/signUp", m.Handler.signUp)

		auth.Group(func(protected chi.Router) {
			protected.Use(middlewares.RequiresSession(m.sRepo.GetByID, m.uRepo.GetByID))
			protected.Post("/signOut", m.Handler.signOut)
			protected.Post("/ping", m.Handler.ping)
			protected.Post("/refresh", m.Handler.refresh)
		})

		auth.Post("/forgot/password/confirm", m.Handler.forgotPasswordConfirm)
	})
}
