package models

import "github.com/juevigrace/diva-server/storage/db"

type SessionStatus int

const (
	SESSION_ACTIVE SessionStatus = iota
	SESSION_EXPIRED
	SESSION_CLOSED
)

func (s SessionStatus) String() string {
	switch s {
	case SESSION_CLOSED:
		return "CLOSED"
	case SESSION_EXPIRED:
		return "EXPIRED"
	case SESSION_ACTIVE:
		return "ACTIVE"
	default:
		return "ACTIVE"
	}
}

func (s SessionStatus) ToDB() db.SessionStatusType {
	switch s {
	case SESSION_CLOSED:
		return db.SessionStatusTypeCLOSED
	case SESSION_EXPIRED:
		return db.SessionStatusTypeEXPIRED
	case SESSION_ACTIVE:
		return db.SessionStatusTypeACTIVE
	default:
		return db.SessionStatusTypeACTIVE
	}
}

func SessionStatusFromString(status string) SessionStatus {
	switch status {
	case "CLOSED":
		return SESSION_CLOSED
	case "EXPIRED":
		return SESSION_EXPIRED
	case "ACTIVE":
		return SESSION_ACTIVE
	default:
		return SESSION_ACTIVE
	}
}

func SessionStatusFromDB(s db.SessionStatusType) SessionStatus {
	switch s {
	case db.SessionStatusTypeCLOSED:
		return SESSION_CLOSED
	case db.SessionStatusTypeEXPIRED:
		return SESSION_EXPIRED
	case db.SessionStatusTypeACTIVE:
		return SESSION_ACTIVE
	default:
		return SESSION_ACTIVE
	}
}
