package eval

import "context"

type _get_interaction_environment struct{}

func (_ *_get_interaction_environment) Type() string {
	return "subroutine.interaction-environment"
}

func (_ *_get_interaction_environment) String() string {
	return "#<subr interaction-environment>"
}

func (_ *_get_interaction_environment) IsList() bool {
	return false
}

func (i *_get_interaction_environment) Equals(sexp SExpression) bool {
	return i.Type() == sexp.Type()
}

func (_ *_get_interaction_environment) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {
	return env.GetGlobalEnv(), nil
}

func NewInteractionEnvironment() SExpression {
	return &_get_interaction_environment{}
}
