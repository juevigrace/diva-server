package postgres

import (
	"errors"
	"fmt"
	"strings"

	"github.com/juevigrace/diva-server/pkg/config"
	"github.com/juevigrace/diva-server/pkg/validator"
)

type PGConf struct {
	Name     string `json:"name" validate:"required"`
	Host     string `json:"host" validate:"required"`
	Port     uint16 `json:"port" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Schema   string `json:"schema"`
	Options  string `json:"options"`
}

func NewPGConf() config.Config {
	var conf *PGConf = new(PGConf)
	conf.LoadDefault()
	return conf
}

func (c *PGConf) Merge(config config.Config) error {
	dc, ok := config.(*PGConf)
	if !ok {
		return errors.New("database: incorrect config type")
	}

	if dc == nil {
		return nil
	}

	if err := dc.Validate(); err != nil {
		return err
	}

	c.Name = dc.Name
	c.Host = dc.Host
	c.Port = dc.Port
	c.Username = dc.Username
	c.Password = dc.Password
	c.Schema = dc.Schema
	c.Options = dc.Options

	return nil
}

func (c *PGConf) LoadFromEnv() {
	c.Name = config.GetEnvOrDefault(DB_NAME_KEY, c.Name)
	c.Host = config.GetEnvOrDefault(DB_HOST_KEY, c.Host)
	c.Port = config.GetEnvOrDefault(DB_PORT_KEY, c.Port)
	c.Username = config.GetEnvOrDefault(DB_USER_KEY, c.Username)
	c.Password = config.GetEnvOrDefault(DB_PASSWORD_KEY, c.Password)
	c.Options = config.GetEnvOrDefault(DB_OPTIONS_KEY, c.Options)
	c.Schema = config.GetEnvOrDefault(DB_SCHEMA_KEY, c.Schema)
}

func (c *PGConf) LoadDefault() {
	c.Name = DB_NAME
	c.Host = DB_HOST
	c.Port = DB_PORT
	c.Username = DB_USERNAME
	c.Password = DB_PASSWORD
	c.Schema = DB_SCHEMA
}

func (c *PGConf) Validate() error {
	return validator.GetInstance().Validate(c)
}

func (c *PGConf) Url() (string, error) {
	schema := c.Schema
	if schema == "" {
		schema = "public"
	}

	var url strings.Builder
	fmt.Fprintf(
		&url,
		"postgres://%s:%s@%s:%d/%s?search_path=%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		schema,
	)

	if c.Options != "" {
		options := strings.SplitSeq(c.Options, ",")
		for opt := range options {
			url.WriteString("&")
			url.WriteString(strings.TrimSpace(opt))
		}
	}

	return url.String(), nil
}


