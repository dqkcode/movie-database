package mongodb

import (
	"time"

	"github.com/dqkcode/movie-database/internal/pkg/config/envconfig"
)

type (
	Config struct {
		Addrs    []string      `envconfig:"MONGODB_ADDRS" default:"127.0.0.1:27017"`
		Database string        `envconfig:"MONGODB_DATABASE" default:"movie"`
		Username string        `envconfig:"MONGODB_USERNAME"`
		Password string        `envconfig:"MONGODB_PASSWORD"`
		Timeout  time.Duration `envconfig:"MONGODB_TIMEOUT" default:"10s"`
	}
)

func LoadConfigFromEnv() Config {
	var conf Config
	envconfig.Load(&conf)
	return conf
}
