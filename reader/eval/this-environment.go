package eval

import "context"

type _this_environment struct{}

func (_ *_this_environment) TypeId() string {
	return "subroutine.this_environment"
}

func (_ *_this_environment) AtomId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_this_environment) String() string {
	return "#<subr this_environment>"
}

func (_ *_this_environment) IsList() bool {
	return false
}

func (i *_this_environment) Equals(sexp SExpression) bool {
	return i.TypeId() == sexp.TypeId()
}

func (_ *_this_environment) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {
	return env, nil
}

func NewThisEnvironment() SExpression {
	return &_this_environment{}
}
