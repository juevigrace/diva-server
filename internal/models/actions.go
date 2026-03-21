package models

type Action int

const (
	ActionUserVerification Action = iota
	ActionPasswordVerification
)

func (a Action) String() string {
	switch a {
	case ActionUserVerification:
		return "USER_VERIFICATION"
	case ActionPasswordVerification:
		return "PASSWORD_VERIFICATION"
	default:
		return "UNKNOWN"
	}
}
