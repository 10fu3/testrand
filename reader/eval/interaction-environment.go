package eval

import "context"

type _get_interaction_environment struct{}

func (_ *_get_interaction_environment) Apply(self *Sexpression, ctx context.Context, env *Sexpression, arguments *Sexpression) (*Sexpression, error) {
	return env._env_globalEnv, nil
}

func NewInteractionEnvironment() *Sexpression {
	return CreateSubroutine("interaction-environment", (&_get_interaction_environment{}).Apply)
}
