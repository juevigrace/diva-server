package models

import "github.com/juevigrace/diva-server/storage/db"

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

func (r Role) ToDB() db.RoleType {
	switch r {
	case ROLE_ADMIN:
		return db.RoleTypeADMIN
	case ROLE_MODERATOR:
		return db.RoleTypeMODERATOR
	case ROLE_USER:
		return db.RoleTypeUSER
	default:
		return db.RoleTypeUSER
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

func RoleFromDB(r db.RoleType) Role {
	switch r {
	case db.RoleTypeADMIN:
		return ROLE_ADMIN
	case db.RoleTypeMODERATOR:
		return ROLE_MODERATOR
	case db.RoleTypeUSER:
		return ROLE_USER
	default:
		return ROLE_USER
	}
}
