package api

import (
	"github.com/dqkcode/movie-database/internal/pkg/db/mongodb"
	"github.com/globalsign/mgo"
)

func GetDBSession() *mgo.Session {

	MongoDBConf := mongodb.LoadConfigFromEnv()
	session, err := mongodb.Dial(MongoDBConf)
	if err != nil {
		panic(err)
	}
	return session
}
