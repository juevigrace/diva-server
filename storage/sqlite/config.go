package sqlite

import (
	"errors"

	"github.com/juevigrace/diva-server/internal/config"
	"github.com/juevigrace/diva-server/pkg/validator"
)

type SQLiteConf struct {
	Path string `json:"path" validate:"required"`
}

func NewSQLiteConf() config.Config {
	var conf *SQLiteConf = new(SQLiteConf)
	conf.LoadDefault()
	return conf
}

func (c *SQLiteConf) Merge(config config.Config) error {
	dc, ok := config.(*SQLiteConf)
	if !ok {
		return errors.New("database: incorrect config type")
	}

	if dc == nil {
		return nil
	}

	if err := dc.Validate(); err != nil {
		return err
	}

	c.Path = dc.Path

	return nil
}

func (c *SQLiteConf) LoadFromEnv() {
	c.Path = config.GetEnvOrDefault(DB_PATH_KEY, c.Path)
}

func (c *SQLiteConf) LoadDefault() {
	c.Path = DB_PATH
}

func (c *SQLiteConf) Validate() error {
	return validator.GetInstance().Validate(c)
}

func (c *SQLiteConf) Url() (string, error) {
	return c.Path, nil
}


