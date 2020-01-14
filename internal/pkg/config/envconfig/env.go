package envconfig

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

var envPrefix = ""

func Load(t interface{}) {
	if err := envconfig.Process(envPrefix, t); err != nil {
		logrus.Errorf("config: unable to load config for %T: %s", t, err)
	}
}

//LoadWithPrefix loads the environment variables with prefix into the provided struct
func LoadWithPrefix(prefix string, t interface{}) {
	if err := envconfig.Process(prefix, t); err != nil {
		logrus.Errorf("config: unable to load config for %T: %s", t, err)
	}
}
