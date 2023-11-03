package eval

import (
	"context"
	"errors"
)

type _apply struct{}

func (_ *_apply) Type() string {
	return "subroutine.apply"
}

func (_ *_apply) String() string {
	return "#<subr apply>"
}

func (_ *_apply) IsList() bool {
	return false
}

func (a *_apply) Equals(sexp SExpression) bool {
	return a.Type() == sexp.Type()
}

func (_ *_apply) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {
	arr, err := ToArray(args)
	if err != nil {
		return nil, err
	}
	if len(arr) != 2 {
		return nil, errors.New("malformed apply")
	}
	car := arr[0]
	cdr := arr[1]
	return car.(Callable).Apply(ctx, env, cdr)
}

func NewApply() SExpression {
	return &_apply{}
}
