package storage

import (
	"errors"
	"fmt"
	"strings"

	"github.com/juevigrace/diva-server/internal/models"
)

type DatabaseConf struct {
	Driver   string `json:"driver"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Schema   string `json:"schema"`
	Options  string `json:"options"`
}

func NewDatabaseConf() models.Config {
	var database *DatabaseConf = new(DatabaseConf)
	database.LoadDefault()
	return database
}

func (c *DatabaseConf) Configure(config models.Config) error {
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

func (c *DatabaseConf) LoadDefault() {
	c.Driver = models.GetEnvOrDefault(DB_DRIVER_KEY, DB_DRIVER)
	c.Name = models.GetEnvOrDefault(DB_NAME_KEY, DB_NAME)
	c.Host = models.GetEnvOrDefault(DB_HOST_KEY, DB_HOST)
	c.Port = models.GetEnvOrDefault(DB_PORT_KEY, DB_PORT)
	c.Username = models.GetEnvOrDefault(DB_USER_KEY, DB_USERNAME)
	c.Password = models.GetEnvOrDefault(DB_PASSWORD_KEY, DB_PASSWORD)
	c.Options = models.GetEnvOrDefault(DB_OPTIONS_KEY, "")
	c.Schema = models.GetEnvOrDefault(DB_SCHEMA_KEY, DB_SCHEMA)
}

func (c *DatabaseConf) Validate() error {
	if c.Driver != "pgx" && c.Driver != "postgres" {
		return errors.New("database: driver must be 'pgx' or 'postgres'")
	}
	if c.Host == "" {
		return errors.New("database: host is required")
	}
	if c.Port == 0 {
		return errors.New("database: port is required")
	}
	if c.Username == "" {
		return errors.New("database: username is required")
	}
	if c.Name == "" {
		return errors.New("database: name is required")
	}
	return nil
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
