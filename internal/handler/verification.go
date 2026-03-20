package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/mail"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/repo"
)

type VerificationHandler struct {
	repo *repo.VerificationRepository
	mail *mail.Client
}

func NewVerificationHandler(repo *repo.VerificationRepository, mail *mail.Client) *VerificationHandler {
	return &VerificationHandler{
		repo: repo,
		mail: mail,
	}
}

func (h *VerificationHandler) Routes(r chi.Router) {
	r.Post("/verify", h.Verify)
}

func (h *VerificationHandler) GenerateAndSend(ctx context.Context, userID uuid.UUID, email string) error {
	u, err := h.repo.Create(ctx, userID)
	if err != nil {
		return fmt.Errorf("create verification: %w", err)
	}

	if err := h.mail.SendVerificationEmail(ctx, email, u); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

func (h *VerificationHandler) Verify(w http.ResponseWriter, r *http.Request) {
	var dto *dtos.EmailTokenDto
	dto, err := middlewares.ValidateBody(dto, r)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	_, err = h.repo.Verify(r.Context(), dto.Token)
	if err != nil {
		var res *responses.APIResponse
		if errors.Is(sql.ErrNoRows, err) {
			res = responses.RespondNotFound(nil, err.Error())
		} else {
			res = responses.RespondBadRequest(nil, err.Error())
		}
		responses.WriteJSON(w, res)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Success"))
}
