package api

import (
	"github.com/dqkcode/movie-database/internal/app/policy"
)

func NewPolicyService() *policy.Service {
	conf := policy.LoadConfigFromEnv()
	return policy.NewEnforcer(conf)
}
