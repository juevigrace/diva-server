package models

import (
	"github.com/juevigrace/diva-server/storage"
)

type Role int

const (
	ROLE_USER Role = iota
	ROLE_MODERATOR
	ROLE_ADMIN
)

func (r Role) String() string {
	switch r {
	case ROLE_ADMIN:
		return "ADMIN"
	case ROLE_MODERATOR:
		return "MODERATOR"
	case ROLE_USER:
		return "USER"
	default:
		return "USER"
	}
}

func RoleFromString(role string) Role {
	switch role {
	case "ADMIN":
		return ROLE_ADMIN
	case "MODERATOR":
		return ROLE_MODERATOR
	case "USER":
		return ROLE_USER
	default:
		return ROLE_USER
	}
}

func (r Role) ToDB() storage.RoleType {
	switch r {
	case ROLE_ADMIN:
		return storage.RoleTypeADMIN
	case ROLE_MODERATOR:
		return storage.RoleTypeMODERATOR
	case ROLE_USER:
		return storage.RoleTypeUSER
	default:
		return storage.RoleTypeUSER
	}
}

func RoleFromDB(r storage.RoleType) Role {
	switch r {
	case storage.RoleTypeADMIN:
		return ROLE_ADMIN
	case storage.RoleTypeMODERATOR:
		return ROLE_MODERATOR
	case storage.RoleTypeUSER:
		return ROLE_USER
	default:
		return ROLE_USER
	}
}
