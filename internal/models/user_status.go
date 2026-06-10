package models

import (
	"github.com/juevigrace/diva-server/storage/db"
)

type UserStatus int

const (
	USER_STATUS_ACTIVE UserStatus = iota
	USER_STATUS_SUSPENDED
	USER_STATUS_INACTIVE
)

func (s UserStatus) String() string {
	switch s {
	case USER_STATUS_ACTIVE:
		return "ACTIVE"
	case USER_STATUS_SUSPENDED:
		return "SUSPENDED"
	case USER_STATUS_INACTIVE:
		return "INACTIVE"
	default:
		return "ACTIVE"
	}
}

func UserStatusFromString(status string) UserStatus {
	switch status {
	case "ACTIVE":
		return USER_STATUS_ACTIVE
	case "SUSPENDED":
		return USER_STATUS_SUSPENDED
	case "INACTIVE":
		return USER_STATUS_INACTIVE
	default:
		return USER_STATUS_ACTIVE
	}
}

func (s UserStatus) ToDB() db.UserStatusType {
	switch s {
	case USER_STATUS_ACTIVE:
		return db.UserStatusTypeACTIVE
	case USER_STATUS_SUSPENDED:
		return db.UserStatusTypeSUSPENDED
	case USER_STATUS_INACTIVE:
		return db.UserStatusTypeINACTIVE
	default:
		return db.UserStatusTypeACTIVE
	}
}

func UserStatusFromDB(s db.UserStatusType) UserStatus {
	switch s {
	case db.UserStatusTypeACTIVE:
		return USER_STATUS_ACTIVE
	case db.UserStatusTypeSUSPENDED:
		return USER_STATUS_SUSPENDED
	case db.UserStatusTypeINACTIVE:
		return USER_STATUS_INACTIVE
	default:
		return USER_STATUS_ACTIVE
	}
}
