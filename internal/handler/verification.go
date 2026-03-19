package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/mail"
	"github.com/juevigrace/diva-server/internal/repo"
	"github.com/juevigrace/diva-server/storage/db"
)

type VerificationHandler struct {
	repo *repo.VerificationRepository
	mail *mail.Client
}

func NewVerificationHandler(queries *db.Queries, mail *mail.Client) *VerificationHandler {
	return &VerificationHandler{
		repo: repo.NewVerificationRepository(queries),
		mail: mail,
	}
}

func (h *VerificationHandler) Routes(r chi.Router) {
	r.Post("/verify", func(w http.ResponseWriter, r *http.Request) {})
}

func (h *VerificationHandler) GenerateAndSend(ctx context.Context, userID *uuid.UUID, email string) error {
	u, err := h.repo.Create(ctx, userID)
	if err != nil {
		return fmt.Errorf("create verification: %w", err)
	}

	if err := h.mail.SendVerificationEmail(ctx, email, u); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

func (h *VerificationHandler) Verify(ctx context.Context, token string) error {
	record, err := h.repo.GetByToken(ctx, token)
	if err != nil {
		return errors.New("invalid token")
	}
	defer func() {
		if record != nil {
			h.repo.DeleteByToken(context.Background(), token)
		}
	}()

	if record.ExpiresAt.Before(time.Now()) {
		return errors.New("token expired")
	}

	return nil
}

func (h *VerificationHandler) DeleteByToken(ctx context.Context, token string) error {
	return h.repo.DeleteByToken(ctx, token)
}
