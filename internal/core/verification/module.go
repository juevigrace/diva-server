package verification

import (
	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/core/user"
	"github.com/juevigrace/diva-server/internal/mail"
	"github.com/juevigrace/diva-server/storage/db"
)

type VerificationModule struct {
	Handler *VerificationHandler
	Repo *VerificationRepo
}

func NewVerificationModule(mail *mail.Client, queries *db.Queries, uModule *user.UserModule) *VerificationModule {
	repo := NewVerificationRepo(mail, queries, uModule.URepo, uModule.UARepo, uModule.UPRepo, uModule.USRepo)
	return &VerificationModule{
		Handler: NewVerificationHandler(repo),
		Repo: repo,
	}
}

func (m *VerificationModule) Routes(r chi.Router) {
	r.Route("/verification", func(v chi.Router) {
		v.Post("/request", m.Handler.requestVerification)
		v.Post("/", m.Handler.verify)
	})
}
