package models

import (
	"github.com/juevigrace/diva-server/storage"
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

func (s UserStatus) ToDB() storage.UserStatusType {
	switch s {
	case USER_STATUS_ACTIVE:
		return storage.UserStatusTypeACTIVE
	case USER_STATUS_SUSPENDED:
		return storage.UserStatusTypeSUSPENDED
	case USER_STATUS_INACTIVE:
		return storage.UserStatusTypeINACTIVE
	default:
		return storage.UserStatusTypeACTIVE
	}
}

func UserStatusFromDB(s storage.UserStatusType) UserStatus {
	switch s {
	case storage.UserStatusTypeACTIVE:
		return USER_STATUS_ACTIVE
	case storage.UserStatusTypeSUSPENDED:
		return USER_STATUS_SUSPENDED
	case storage.UserStatusTypeINACTIVE:
		return USER_STATUS_INACTIVE
	default:
		return USER_STATUS_ACTIVE
	}
}
