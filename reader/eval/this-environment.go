package eval

import "context"

func _subr_this_environment_Apply(self *Sexpression, ctx context.Context, env *Sexpression, arguments *Sexpression) (*Sexpression, error) {
	return env, nil
}

func NewThisEnvironment() *Sexpression {
	return CreateSubroutine("this-environment", _subr_this_environment_Apply)
}
