package storage

import (
	"errors"
	"fmt"
	"strings"

	"github.com/juevigrace/diva-server/internal/config"
	"github.com/juevigrace/diva-server/pkg/validator"
)

type DatabaseConf struct {
	Driver   string `json:"driver" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Host     string `json:"host" validate:"required"`
	Port     uint16 `json:"port" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Schema   string `json:"schema"`
	Options  string `json:"options"`
}

func NewDatabaseConf() *DatabaseConf {
	var database *DatabaseConf = new(DatabaseConf)
	database.LoadDefault()
	return database
}

func (c *DatabaseConf) Merge(config config.Config) error {
	dc, ok := config.(*DatabaseConf)
	if !ok {
		return errors.New("database: incorrect config type")
	}

	if dc == nil {
		return nil
	}

	if err := dc.Validate(); err != nil {
		return err
	}

	c.Driver = dc.Driver
	c.Name = dc.Name
	c.Host = dc.Host
	c.Port = dc.Port
	c.Username = dc.Username
	c.Password = dc.Password
	c.Schema = dc.Schema
	c.Options = dc.Options

	return nil
}

func (c *DatabaseConf) LoadFromFile(path string) error {
	return nil
}

func (c *DatabaseConf) LoadFromEnv() {
	c.Driver = config.GetEnvOrDefault(DB_DRIVER_KEY, c.Driver)
	c.Name = config.GetEnvOrDefault(DB_NAME_KEY, c.Name)
	c.Host = config.GetEnvOrDefault(DB_HOST_KEY, c.Host)
	c.Port = config.GetEnvOrDefault(DB_PORT_KEY, c.Port)
	c.Username = config.GetEnvOrDefault(DB_USER_KEY, c.Username)
	c.Password = config.GetEnvOrDefault(DB_PASSWORD_KEY, c.Password)
	c.Options = config.GetEnvOrDefault(DB_OPTIONS_KEY, c.Options)
	c.Schema = config.GetEnvOrDefault(DB_SCHEMA_KEY, c.Schema)
}

func (c *DatabaseConf) LoadDefault() {
	c.Driver = DB_DRIVER
	c.Name = DB_NAME
	c.Host = DB_HOST
	c.Port = DB_PORT
	c.Username = DB_USERNAME
	c.Password = DB_PASSWORD
	c.Schema = DB_SCHEMA
}

func (c *DatabaseConf) SaveToFile(path string) error {
	return nil
}

func (c *DatabaseConf) Validate() error {
	return validator.GetInstance().Validate(c)
}

func (c *DatabaseConf) Url() (string, error) {
	schema := c.Schema
	if schema == "" {
		schema = "public"
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?search_path=%s",
		c.Username, c.Password, c.Host, c.Port, c.Name, schema)

	if c.Options != "" {
		options := strings.Split(c.Options, ",")
		for _, opt := range options {
			url += "&" + strings.TrimSpace(opt)
		}
	}

	return url, nil
}
