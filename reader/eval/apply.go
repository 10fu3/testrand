package eval

import (
	"context"
	"errors"
)

type _apply struct{}

func (_ *_apply) TypeId() string {
	return "subroutine.apply"
}

func (_ *_apply) AtomId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_apply) String() string {
	return "#<subr apply>"
}

func (_ *_apply) IsList() bool {
	return false
}

func (a *_apply) Equals(sexp SExpression) bool {
	return a.TypeId() == sexp.TypeId()
}

func (_ *_apply) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {
	if argsLength != 2 {
		return nil, errors.New("malformed apply")
	}
	car := args[0]
	cdr := args[1]

	cdrArr, size, err := ToArray(cdr)

	if err != nil {
		return nil, err
	}

	return car.(Callable).Apply(ctx, env, cdrArr, size)
}

func NewApply() SExpression {
	return &_apply{}
}
