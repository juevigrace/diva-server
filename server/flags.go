package server

import "flag"

const (
	CONFIG_PATH = "./config.json"
	USE_ENV     = false
)

type ServerFlags struct {
	ConfigPath string
	UsesEnv    bool
}

func NewServerFlags() *ServerFlags {
	flags := new(ServerFlags)
	flags.Parse()
	return flags
}

func (f *ServerFlags) Parse() {
	var config string
	var env bool
	flag.BoolVar(&env, "env", USE_ENV, "Tells the server to use environment variables")
	flag.BoolVar(&env, "e", USE_ENV, "Tells the server to use environment variables")
	flag.StringVar(&config, "config", CONFIG_PATH, "Tells the server to load the config file")
	flag.StringVar(&config, "c", CONFIG_PATH, "Tells the server to load the config file")
	flag.Parse()

	f.ConfigPath = config
	f.UsesEnv = env
}
