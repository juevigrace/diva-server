package models

type SessionStatus int

const (
	SESSION_ACTIVE SessionStatus = iota
	SESSION_EXPIRED
	SESSION_CLOSED
)

func (s SessionStatus) String() string {
	switch s {
	case SESSION_CLOSED:
		return "closed"
	case SESSION_EXPIRED:
		return "expired"
	case SESSION_ACTIVE:
		return "active"
	default:
		return "active"
	}
}

func SessionStatusFromString(status string) SessionStatus {
	switch status {
	case "closed":
		return SESSION_CLOSED
	case "expired":
		return SESSION_EXPIRED
	case "active":
		return SESSION_ACTIVE
	default:
		return SESSION_ACTIVE
	}
}
