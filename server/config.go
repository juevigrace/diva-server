package server

import (
	"errors"

	"github.com/juevigrace/diva-server/internal/config"
	"github.com/juevigrace/diva-server/pkg/validator"
)

type ServerConfig struct {
	Version         string
	Port            uint16
	Domain          string
	Env             config.Env
	Debug           bool
	UploadsDir      string
	JWTSecret       string
	ResendAPIKey    string
	ResendFromEmail string
	RootUsername    string
	RootPassword    string
	RootEmail       string
}

func NewServerConfig() config.Config {
	var c *ServerConfig = new(ServerConfig)
	c.LoadDefault()
	return c
}

func (c *ServerConfig) Merge(from config.Config) error {
	sc, ok := from.(*ServerConfig)
	if !ok {
		return errors.New("server: incorrect config type")
	}

	if sc == nil {
		return nil
	}

	if err := sc.Validate(); err != nil {
		return err
	}

	c.Port = sc.Port
	c.Domain = sc.Domain
	c.Env = sc.Env
	c.Debug = sc.Debug
	c.UploadsDir = sc.UploadsDir
	c.JWTSecret = sc.JWTSecret
	c.ResendAPIKey = sc.ResendAPIKey
	c.ResendFromEmail = sc.ResendFromEmail
	c.RootUsername = sc.RootUsername
	c.RootPassword = sc.RootPassword
	c.RootEmail = sc.RootEmail

	return nil
}

func (c *ServerConfig) LoadFromEnv() {
	c.Port = config.GetEnvOrDefault(SERVER_PORT_KEY, c.Port)
	c.Domain = config.GetEnvOrDefault(SERVER_DOMAIN_KEY, c.Domain)
	c.Debug = config.GetEnvOrDefault(SERVER_DEBUG_KEY, c.Debug)

	envStr := config.GetEnvOrDefault(SERVER_ENV_KEY, c.Env.String())
	c.Env = config.StringToEnv(envStr)

	c.JWTSecret = config.GetEnvOrDefault(JWT_SECRET_KEY, c.JWTSecret)
	c.ResendAPIKey = config.GetEnvOrDefault(RESEND_API_KEY, c.ResendAPIKey)
	c.ResendFromEmail = config.GetEnvOrDefault(RESEND_FROM_EMAIL, c.ResendFromEmail)
	c.RootUsername = config.GetEnvOrDefault(ROOT_USERNAME_KEY, c.RootUsername)
	c.RootPassword = config.GetEnvOrDefault(ROOT_PASSWORD_KEY, c.RootPassword)
	c.RootEmail = config.GetEnvOrDefault(ROOT_EMAIL_KEY, c.RootEmail)
	c.UploadsDir = config.GetEnvOrDefault(UPLOADS_DIR_KEY, c.UploadsDir)
}

func (c *ServerConfig) LoadDefault() {
	c.Version = VERSION
	c.Port = SERVER_PORT
	c.Domain = SERVER_DOMAIN
	c.Env = SERVER_ENV
	c.Debug = SERVER_DEBUG
	c.JWTSecret = JWT_SECRET_DEFAULT
	c.ResendAPIKey = RESEND_API_KEY_DEFAULT
	c.ResendFromEmail = RESEND_FROM_EMAIL_DEF
	c.RootUsername = ROOT_USERNAME
	c.RootPassword = ROOT_PASSWORD
	c.RootEmail = ROOT_EMAIL
	c.UploadsDir = UPLOADS_DIR
}

func (c *ServerConfig) Validate() error {
	return validator.GetInstance().Validate(c)
}
