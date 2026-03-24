package models

type Action int

const (
	ActionUserVerification Action = iota
)

func (a Action) String() string {
	switch a {
	case ActionUserVerification:
		return "USER_VERIFICATION"
	default:
		return "UNKNOWN"
	}
}

func ActionFromString(s string) Action {
	switch s {
	case "USER_VERIFICATION":
		return ActionUserVerification
	default:
		return -1
	}
}
