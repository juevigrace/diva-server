package models

import "github.com/juevigrace/diva-server/storage"

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

func (s SessionStatus) ToDB() storage.SessionStatusType {
	switch s {
	case SESSION_CLOSED:
		return storage.SessionStatusTypeCLOSED
	case SESSION_EXPIRED:
		return storage.SessionStatusTypeEXPIRED
	case SESSION_ACTIVE:
		return storage.SessionStatusTypeACTIVE
	default:
		return storage.SessionStatusTypeACTIVE
	}
}

func SessionStatusFromDB(s storage.SessionStatusType) SessionStatus {
	switch s {
	case storage.SessionStatusTypeCLOSED:
		return SESSION_CLOSED
	case storage.SessionStatusTypeEXPIRED:
		return SESSION_EXPIRED
	case storage.SessionStatusTypeACTIVE:
		return SESSION_ACTIVE
	default:
		return SESSION_ACTIVE
	}
}
