package storage

const (
	DB_DRIVER   string = "pgx"
	DB_NAME     string = "diva"
	DB_HOST     string = "localhost"
	DB_PORT     uint16 = 5432
	DB_USERNAME string = "diva_user"
	DB_PASSWORD string = "secret_password"
	DB_SCHEMA   string = "public"
)

const (
	DB_DRIVER_KEY   string = "DB_DRIVER"
	DB_HOST_KEY     string = "DB_HOST"
	DB_PORT_KEY     string = "DB_PORT"
	DB_USER_KEY     string = "DB_USER"
	DB_PASSWORD_KEY string = "DB_PASSWORD"
	DB_NAME_KEY     string = "DB_NAME"
	DB_SCHEMA_KEY   string = "DB_SCHEMA"
	DB_OPTIONS_KEY  string = "DB_OPTIONS"
)
