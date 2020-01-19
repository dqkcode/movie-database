package server

import (
	"time"

	"github.com/dqkcode/movie-database/internal/pkg/config/envconfig"
)

type (
	Config struct {
		Address      string        `envconfig:"HTTP_ADDRESS" default:""`
		Port         int           `envconfig:"HTTP_PORT" default:"8080"`
		ReadTimeOut  time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"5m"`
		WriteTimeOut time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"5m"`
	}
)

func LoadConfigFromEnv() Config {
	var conf Config
	envconfig.Load(&conf)
	return conf
}
