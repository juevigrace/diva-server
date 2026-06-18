package server

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/juevigrace/diva-server/internal/config"
	"github.com/juevigrace/diva-server/pkg/validator"
	"github.com/juevigrace/diva-server/storage"
)

type ServerConfig struct {
	Flags           *ServerFlags          `json:"-"`
	Version         string                `json:"-"`
	Port            uint16                `json:"port" validate:"required"`
	Domain          string                `json:"domain" validate:"required"`
	Env             config.Env            `json:"env" validate:"max=1"`
	Debug           bool                  `json:"debug"`
	JWTSecret       string                `json:"jwt_secret" validate:"required"`
	ResendAPIKey    string                `json:"resend_api_key" validate:"required"`
	ResendFromEmail string                `json:"resend_from_email" validate:"required"`
	RootUsername    string                `json:"root_username" validate:"required"`
	RootPassword    string                `json:"root_password" validate:"required"`
	RootEmail       string                `json:"root_email" validate:"required"`
	UploadsDir      string                `json:"uploads_dir" validate:"required"`
	Database        *storage.DatabaseConf `json:"database" validate:"required"`
}

func NewServerConfig(flags *ServerFlags) config.Config {
	var config *ServerConfig = new(ServerConfig)
	config.Flags = flags
	config.LoadDefault()
	return config
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
	c.UploadsDir = sc.UploadsDir

	if err := c.Database.Merge(sc.Database); err != nil {
		return err
	}

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

	c.Database.LoadFromEnv()
}

func (c *ServerConfig) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, c)
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

	c.Database = storage.NewDatabaseConf()
}

func (c *ServerConfig) SaveToFile(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (c *ServerConfig) Validate() error {
	return validator.GetInstance().Validate(c)
}
