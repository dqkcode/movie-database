package policy

import (
	"github.com/casbin/casbin"
	"github.com/dqkcode/movie-database/internal/app/auth"
	"github.com/dqkcode/movie-database/internal/app/types"
	"github.com/dqkcode/movie-database/internal/pkg/config/envconfig"
	"golang.org/x/net/context"
)

type (
	Config struct {
		Model  string `envconfig:"CONFIG_MODEL" default:"configs/acl_model.conf"`
		Policy string `envconfig:"CONFIG_POLICY" default:"configs/acl_policy.csv"`
	}
	Service struct {
		enforcer *casbin.Enforcer
	}
)

func LoadConfigFromEnv() Config {
	var conf Config
	envconfig.Load(&conf)
	return conf
}

func NewEnforcer(c Config) *Service {
	return &Service{
		enforcer: casbin.NewEnforcer(c.Model, c.Policy),
	}
}

func (s *Service) IsAllowed(sub, obj, act string) bool {
	ok, err := s.enforcer.EnforceSafe(sub, obj, act)
	return err == nil && ok
}

func (s *Service) Validate(ctx context.Context, obj, act string) bool {

	sub := auth.GetRoleFromContext(ctx)
	if sub == types.RoleAdmin {
		return true
	}
	return s.IsAllowed((string)(sub), obj, act)
}
