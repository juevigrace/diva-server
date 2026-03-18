package server

import (
	"errors"

	"github.com/juevigrace/diva-server/internal/models"
)

type ServerConfig struct {
	Version         string           `json:"-"`
	Port            uint16           `json:"port"`
	Domain          string           `json:"domain"`
	Env             models.ServerEnv `json:"env"`
	Debug           bool             `json:"debug"`
	JWTSecret       string           `json:"jwt_secret"`
	ResendAPIKey    string           `json:"resend_api_key"`
	ResendFromEmail string           `json:"resend_from_email"`
	RootUsername    string           `json:"root_username"`
	RootPassword    string           `json:"root_password"`
	RootEmail       string           `json:"root_email"`
}

func NewServerConfig() models.Config {
	var config *ServerConfig = new(ServerConfig)
	config.LoadDefault()
	return config
}

func (c *ServerConfig) Configure(config models.Config) error {
	sc, ok := config.(*ServerConfig)
	if !ok {
		return errors.New("server: incorrect config type")
	}

	if sc == nil {
		return nil
	}

	if err := sc.Validate(); err != nil {
		return err
	}

	c.Version = sc.Version
	c.Port = sc.Port
	c.Domain = sc.Domain
	c.Env = sc.Env
	c.Debug = sc.Debug
	c.JWTSecret = sc.JWTSecret
	c.ResendAPIKey = sc.ResendAPIKey
	c.ResendFromEmail = sc.ResendFromEmail
	c.RootUsername = sc.RootUsername
	c.RootPassword = sc.RootPassword
	c.RootEmail = sc.RootEmail

	return nil
}

func (c *ServerConfig) LoadDefault() {
	c.Version = VERSION
	c.Port = models.GetEnvOrDefault(SERVER_PORT_KEY, SERVER_PORT)
	c.Domain = models.GetEnvOrDefault(SERVER_DOMAIN_KEY, SERVER_DOMAIN)
	c.Debug = models.GetEnvOrDefault(SERVER_DEBUG_KEY, SERVER_DEBUG)

	envStr := models.GetEnvOrDefault(SERVER_ENV_KEY, SERVER_ENV.String())
	c.Env = models.StringToServerEnv(envStr)

	c.JWTSecret = models.GetEnvOrDefault(JWT_SECRET_KEY, JWT_SECRET_DEFAULT)
	c.ResendAPIKey = models.GetEnvOrDefault(RESEND_API_KEY, RESEND_API_KEY_DEFAULT)
	c.ResendFromEmail = models.GetEnvOrDefault(RESEND_FROM_EMAIL, RESEND_FROM_EMAIL_DEF)
	c.RootUsername = models.GetEnvOrDefault(ROOT_USERNAME_KEY, ROOT_USERNAME)
	c.RootPassword = models.GetEnvOrDefault(ROOT_PASSWORD_KEY, ROOT_PASSWORD)
	c.RootEmail = models.GetEnvOrDefault(ROOT_EMAIL_KEY, ROOT_EMAIL)
}

func (c *ServerConfig) Validate() error {
	if c.Port == 0 {
		return errors.New("server: port is required")
	}
	if c.Domain == "" {
		return errors.New("server: domain is required")
	}
	if c.JWTSecret == "" {
		return errors.New("server: jwt secret is required")
	}
	if c.RootUsername == "" {
		return errors.New("server: root username is required")
	}
	if c.RootPassword == "" {
		return errors.New("server: root password is required")
	}
	if c.RootEmail == "" {
		return errors.New("server: root email is required")
	}
	return nil
}
