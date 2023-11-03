package eval

import "context"

type _this_environment struct{}

func (_ *_this_environment) Type() string {
	return "subroutine.this_environment"
}

func (_ *_this_environment) String() string {
	return "#<subr this_environment>"
}

func (_ *_this_environment) IsList() bool {
	return false
}

func (i *_this_environment) Equals(sexp SExpression) bool {
	return i.Type() == sexp.Type()
}

func (_ *_this_environment) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {
	return env, nil
}

func NewThisEnvironment() SExpression {
	return &_this_environment{}
}
