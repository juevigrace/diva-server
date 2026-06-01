package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/storage/db"
)

var (
	ErrTokenExpired      = errors.New("token expired")
	ErrTokenInvalid      = errors.New("token invalid")
	ErrActionNotFound    = errors.New("action not found")
	ErrActionNotVerified = errors.New("action not verified")
)

type UserAction struct {
	ID     uuid.UUID
	Name   Action
	UserID uuid.UUID
}

type UserActionVerification struct {
	Action    UserAction
	Token     string
	ExpiresAt time.Time
	UsedAt    time.Time
	Verified  bool
}

type Action int

const (
	ActionUserVerification Action = iota
	ActionPasswordReset
)

func (a Action) String() string {
	switch a {
	case ActionUserVerification:
		return "USER_VERIFICATION"
	case ActionPasswordReset:
		return "PASSWORD_RESET"
	default:
		return ""
	}
}

func ActionFromString(s string) Action {
	switch s {
	case "USER_VERIFICATION":
		return ActionUserVerification
	case "PASSWORD_RESET":
		return ActionPasswordReset
	default:
		return -1
	}
}

func (ua *UserAction) Response() *responses.UserActionResponse {
	return &responses.UserActionResponse{
		ID:         ua.ID.String(),
		ActionName: ua.Name.String(),
	}
}

func (ua *UserAction) DBCreate() *db.CreateUserActionParams {
	return &db.CreateUserActionParams{
		ID:     UUIDPtrToDB(&ua.ID),
		Name:   ua.Name.String(),
		UserID: UUIDPtrToDB(&ua.UserID),
	}
}

func (uv *UserActionVerification) DBCreate() *db.CreateUserVerificationParams {
	return &db.CreateUserVerificationParams{
		ActionID: pgtype.UUID{
			Bytes: uv.Action.ID,
			Valid: true,
		},
		Token: uv.Token,
		ExpiresAt: pgtype.Timestamptz{
			Time:  uv.ExpiresAt,
			Valid: true,
		},
	}
}

func UserActionFromDB(row *db.DivaAction) *UserAction {
	return &UserAction{
		ID:     DBUUIDToUUID(row.ID),
		Name:   ActionFromString(row.Name),
		UserID: DBUUIDToUUID(row.UserID),
	}
}
