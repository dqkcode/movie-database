package mongodb

import (
	"time"

	"github.com/dqkcode/movie-database/internal/pkg/config/envconfig"
	"github.com/globalsign/mgo"
	"github.com/sirupsen/logrus"
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

func LoadConfigFromEnv() *Config {
	var conf Config
	envconfig.Load(&conf)
	return &conf
}

func Dial(config *Config) (*mgo.Session, error) {
	logrus.Infof("dialing to MongoDB at: %v, database: %v", config.Addrs, config.Database)
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    config.Addrs,
		Database: config.Database,
		Password: config.Password,
		Username: config.Username,
		Timeout:  config.Timeout,
	})
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	logrus.Infof("successfully dialing to MongoDB at %v", config.Addrs)

	return session, nil
}
func InitDBSession() *mgo.Session {

	MongoDBConf := LoadConfigFromEnv()
	session, err := Dial(MongoDBConf)
	if err != nil {
		logrus.Errorf("Can not dial to MongoDB addrs: %v", MongoDBConf.Addrs)
		panic(err)
	}
	return session
}
