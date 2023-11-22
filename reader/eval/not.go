package eval

import (
	"context"
	"errors"
)

type _is_not struct {
}

func (_ *_is_not) TypeId() string {
	return "subroutine.is_equals"
}

func (_ *_is_not) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_is_not) String() string {
	return "#<subr eq?>"
}

func (_ *_is_not) IsList() bool {
	return false
}

func (i *_is_not) Equals(sexp SExpression) bool {
	return i.TypeId() == sexp.TypeId()
}

func (_ *_is_not) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength != 1 {
		return nil, errors.New("malformed not")
	}

	first := args[0]

	if "bool" != first.TypeId() {
		return first, nil
	}

	return NewBool(!first.(Bool).GetValue()), nil
}

func NewIsNot() SExpression {
	return &_is_not{}
}
