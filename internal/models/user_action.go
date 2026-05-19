package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models/responses"
)

var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenInvalid = errors.New("token invalid")
)

type UserAction struct {
	ID     uuid.UUID
	Action Action
}

type UserActionVerification struct {
	Action    UserAction
	Token     string
	ExpiresAt time.Time
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
		return "UNKNOWN"
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

func (ua *UserAction) Response(id *uuid.UUID) *responses.UserActionResponse {
	var userID *string = new(string)
	if id != nil {
		*userID = id.String()
	}

	return &responses.UserActionResponse{
		UserID:     userID,
		ID:         ua.ID.String(),
		ActionName: ua.Action.String(),
	}
}
