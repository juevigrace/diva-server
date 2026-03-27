package models

type Action int

const (
	ActionUserVerification Action = iota
	ActionPasswordReset
)

func (a Action) String() string {
	switch a {
	case ActionUserVerification:
		return "USER_VERIFICATION"
	case ActionPasswordReset:
		return "PASSWORD_RESET"
	default:
		return "UNKNOWN"
	}
}

func ActionFromString(s string) Action {
	switch s {
	case "USER_VERIFICATION":
		return ActionUserVerification
	case "PASSWORD_RESET":
		return ActionPasswordReset
	default:
		return -1
	}
}
