package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/storage"
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
	UsedAt    *time.Time
	Verified  bool
}

type Action int

const (
	ActionUserVerification Action = iota
	ActionPasswordUpdate
	ActionEmailUpdate
	ActionUsernameUpdate
	ActionPhoneUpdate
	ActionUserRestore
)

func (a Action) String() string {
	switch a {
	case ActionUserVerification:
		return "USER_VERIFICATION"
	case ActionPasswordUpdate:
		return "PASSWORD_RESET"
	case ActionEmailUpdate:
		return "EMAIL_UPDATE"
	case ActionUsernameUpdate:
		return "USERNAME_UPDATE"
	case ActionPhoneUpdate:
		return "PHONE_UPDATE"
	case ActionUserRestore:
		return "USER_RESTORE"
	default:
		return ""
	}
}

func ActionFromString(s string) Action {
	switch s {
	case "USER_VERIFICATION":
		return ActionUserVerification
	case "PASSWORD_RESET":
		return ActionPasswordUpdate
	case "EMAIL_UPDATE":
		return ActionEmailUpdate
	case "USERNAME_UPDATE":
		return ActionUsernameUpdate
	case "PHONE_UPDATE":
		return ActionPhoneUpdate
	case "USER_RESTORE":
		return ActionUserRestore
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

func (ua *UserAction) DBCreate() *storage.CreateUserActionParams {
	return &storage.CreateUserActionParams{
		ID:     ua.ID,
		Name:   ua.Name.String(),
		UserID: ua.UserID,
	}
}

func (uv *UserActionVerification) DBCreate() *storage.CreateUserVerificationParams {
	return &storage.CreateUserVerificationParams{
		ActionID:  uv.Action.ID,
		Token:     uv.Token,
		ExpiresAt: uv.ExpiresAt.UnixMilli()}
}

func UserActionFromDB(row *storage.DivaAction) *UserAction {
	return &UserAction{
		ID:     row.ID,
		Name:   ActionFromString(row.Name),
		UserID: row.UserID,
	}
}
