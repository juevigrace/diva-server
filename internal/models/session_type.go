package models

import "github.com/juevigrace/diva-server/storage"

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

func (st *SessionType) ToDB() storage.SessionType {
	switch *st {
	case SESSION_NORMAL:
		return storage.SessionTypeNORMAL
	case SESSION_TEMPORAL:
		return storage.SessionTypeTEMPORAL
	default:
		return storage.SessionTypeNORMAL
	}
}

func SessionTypeFromDB(t storage.SessionType) SessionType {
	switch t {
	case storage.SessionTypeNORMAL:
		return SESSION_NORMAL
	case storage.SessionTypeTEMPORAL:
		return SESSION_TEMPORAL
	default:
		return SESSION_NORMAL
	}
}
