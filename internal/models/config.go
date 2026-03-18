package models

import (
	"os"
	"strconv"
)

type Config interface {
	Configure(config Config) error
	LoadDefault()
	Validate() error
}

func GetEnvOrDefault[T any](name string, defaultValue T) T {
	env, exists := os.LookupEnv(name)
	if !exists {
		return defaultValue
	}

	var result T
	var err error

	switch any(result).(type) {
	case uint16:
		var val uint64
		val, err = strconv.ParseUint(env, 10, 16)
		if err == nil {
			result = any(uint16(val)).(T)
		}
	case bool:
		var val bool
		val, err = strconv.ParseBool(env)
		if err == nil {
			result = any(val).(T)
		}
	case string:
		result = any(env).(T)
	}

	if err != nil {
		return defaultValue
	}
	return result
}
