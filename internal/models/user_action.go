package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/storage/db"
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

func (ua *UserAction) DBCreate() *db.CreateUserActionParams {
	return &db.CreateUserActionParams{
		ID:     UUIDPtrToDB(&ua.ID),
		Name:   ua.Name.String(),
		UserID: UUIDPtrToDB(&ua.UserID),
	}
}

func (uv *UserActionVerification) DBCreate() *db.CreateUserVerificationParams {
	exp := uv.ExpiresAt.UnixMilli()
	return &db.CreateUserVerificationParams{
		ActionID:  UUIDPtrToDB(&uv.Action.ID),
		Token:     uv.Token,
		ExpiresAt: IntPtrToDBTime(&exp)}
}

func UserActionFromDB(row *db.DivaAction) *UserAction {
	return &UserAction{
		ID:     DBUUIDToUUID(row.ID),
		Name:   ActionFromString(row.Name),
		UserID: DBUUIDToUUID(row.UserID),
	}
}
