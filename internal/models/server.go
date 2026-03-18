package models

type ServerEnv int

const (
	DEVELOPMENT ServerEnv = iota
	PRODUCTION
)

func (s ServerEnv) String() string {
	switch s {
	case DEVELOPMENT:
		return "DEVELOPMENT"
	case PRODUCTION:
		return "PRODUCTION"
	default:
		return "DEVELOPMENT"
	}
}

func StringToServerEnv(env string) ServerEnv {
	switch env {
	case "DEVELOPMENT", "DEV", "dev":
		return DEVELOPMENT
	case "PRODUCTION", "PROD", "prod":
		return PRODUCTION
	default:
		return DEVELOPMENT
	}
}

type ServerStatus int

const (
	Inactive ServerStatus = iota
	Active
	Maintenance
)

func (s ServerStatus) String() string {
	switch s {
	case Inactive:
		return "INACTIVE"
	case Active:
		return "ACTIVE"
	case Maintenance:
		return "MAINTENANCE"
	default:
		return "INACTIVE"
	}
}
