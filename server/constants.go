package server

import "github.com/juevigrace/diva-server/internal/models"

const (
	VERSION string = "1.0.0"
)

const (
	SERVER_PORT            uint16           = 5000
	SERVER_DOMAIN          string           = "localhost"
	SERVER_ENV             models.ServerEnv = models.DEVELOPMENT
	SERVER_DEBUG           bool             = true
	JWT_SECRET_DEFAULT     string           = "super secret"
	RESEND_API_KEY_DEFAULT string           = ""
	RESEND_FROM_EMAIL_DEF  string           = "onboarding@resend.dev"
	ROOT_USERNAME          string           = "user"
	ROOT_PASSWORD          string           = "1234"
	ROOT_EMAIL             string           = "user@gmail.com"
)

const (
	SERVER_PORT_KEY   string = "PORT"
	SERVER_DOMAIN_KEY string = "DOMAIN"
	SERVER_ENV_KEY    string = "SERVER_ENV"
	SERVER_DEBUG_KEY  string = "DEBUG"
	JWT_SECRET_KEY    string = "JWT_SECRET"
	RESEND_API_KEY    string = "RESEND_API_KEY"
	RESEND_FROM_EMAIL string = "RESEND_FROM_EMAIL"
	ROOT_USERNAME_KEY string = "ROOT_USERNAME"
	ROOT_PASSWORD_KEY string = "ROOT_PASSWORD"
	ROOT_EMAIL_KEY    string = "ROOT_EMAIL"
)
