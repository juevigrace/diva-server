package models

import "github.com/juevigrace/diva-server/storage/db"

type SessionType int

const (
	SESSION_NORMAL SessionType = iota
	SESSION_TEMPORAL
)

func (st *SessionType) String() string {
	switch *st {
	case SESSION_NORMAL:
		return "NORMAL"
	case SESSION_TEMPORAL:
		return "TEMPORAL"
	default:
		return "NORMAL"
	}
}

func SessionTypeFromString(s string) SessionType {
	switch s {
	case "NORMAL":
		return SESSION_NORMAL
	case "TEMPORAL":
		return SESSION_TEMPORAL
	default:
		return SESSION_NORMAL
	}
}

func (st *SessionType) ToDB() db.SessionType {
	switch *st {
	case SESSION_NORMAL:
		return db.SessionTypeNORMAL
	case SESSION_TEMPORAL:
		return db.SessionTypeTEMPORAL
	default:
		return db.SessionTypeNORMAL
	}
}

func SessionTypeFromDB(t db.SessionType) SessionType {
	switch t {
	case db.SessionTypeNORMAL:
		return SESSION_NORMAL
	case db.SessionTypeTEMPORAL:
		return SESSION_TEMPORAL
	default:
		return SESSION_NORMAL
	}
}
