package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/mail"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/repo"
	"github.com/juevigrace/diva-server/internal/util"
	"github.com/juevigrace/diva-server/storage/db"
)

type VerificationService struct {
	repo *repo.VerificationRepository
	mail *mail.Client
}

func NewVerificationService(repo *repo.VerificationRepository, mail *mail.Client) *VerificationService {
	return &VerificationService{
		repo: repo,
		mail: mail,
	}
}

func (s *VerificationService) GenerateAndSend(ctx context.Context, userID uuid.UUID, email string) error {
	token, err := util.GenerateOTPCode()
	if err != nil {
		return err
	}

	params := &db.CreateVerificationParams{
		UserID:    pgtype.UUID{Bytes: userID, Valid: true},
		Token:     token,
		CreatedAt: pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().UTC().Add(15 * time.Minute), Valid: true},
	}

	if err := s.repo.Create(ctx, params); err != nil {
		return err
	}

	verification, err := s.repo.GetByUser(ctx, userID)
	if err != nil {
		return err
	}

	go func() {
		if err := s.mail.SendVerificationEmail(ctx, email, verification); err != nil {
			log.Println(err)
			return
		}
	}()

	return nil
}

func (s *VerificationService) Verify(ctx context.Context, token string) (*models.UserVerification, error) {
	record, err := s.repo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if record.ExpiresAt.Before(time.Now().UTC()) {
		return nil, errors.New("token expired")
	}

	defer func() {
		err = s.repo.DeleteByToken(ctx, token)
		if err != nil {
			log.Println(err)
		}
	}()

	return record, nil
}
