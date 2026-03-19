package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/repo"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserHandler struct {
	repo         *repo.UserRepository
	verification *VerificationHandler
}

func NewUserHandler(queries *db.Queries, verification *VerificationHandler) *UserHandler {
	return &UserHandler{
		repo:         repo.NewUserRepository(queries),
		verification: verification,
	}
}

func (h *UserHandler) Routes(r chi.Router) {
	r.Route("/user", func(user chi.Router) {
		user.Route("/verify", func(verify chi.Router) {
			// verify.Use(middlewares.SessionMiddleware(h.))
			verify.Post("/email", h.verifyUserEmail)
		})
	})
}

func (h *UserHandler) verifyUserEmail(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "authenticate first"))
		return
	}

	var body *dtos.EmailTokenDto
	body, err := middlewares.ValidateBody(body, r)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	err = h.verification.Verify(r.Context(), body.Token)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.repo.VerifyUser(context.Background(), session.User); err != nil {
		responses.WriteJSON(w, responses.RespondInternalServerError(nil, "failed to verify user"))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "user verified"))
}
