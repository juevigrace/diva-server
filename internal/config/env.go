package config

type Env int

const (
	DEVELOPMENT Env = iota
	PRODUCTION
)

func (s Env) String() string {
	switch s {
	case DEVELOPMENT:
		return "DEVELOPMENT"
	case PRODUCTION:
		return "PRODUCTION"
	default:
		return "DEVELOPMENT"
	}
}

func StringToEnv(env string) Env {
	switch env {
	case "DEVELOPMENT", "DEV", "dev":
		return DEVELOPMENT
	case "PRODUCTION", "PROD", "prod":
		return PRODUCTION
	default:
		return DEVELOPMENT
	}
}
