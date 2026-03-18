package models

type Role int

const (
	ROLE_USER Role = iota
	ROLE_MODERATOR
	USER_ADMIN
)

func (r *Role) String() string {
	switch *r {
	case USER_ADMIN:
		return "admin"
	case ROLE_MODERATOR:
		return "moderator"
	case ROLE_USER:
		return "user"
	default:
		return "user"
	}
}

func RoleFromString(role string) Role {
	switch role {
	case "admin":
		return USER_ADMIN
	case "moderator":
		return ROLE_MODERATOR
	case "user":
		return ROLE_USER
	default:
		return ROLE_USER
	}
}
