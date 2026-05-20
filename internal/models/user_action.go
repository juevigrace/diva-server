package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models/responses"
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
