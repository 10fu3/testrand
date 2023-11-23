package eval

import (
	"context"
	"errors"
)

type _is_equals struct {
}

func (_ *_is_equals) TypeId() string {
	return "subroutine.is_equals"
}

func (_ *_is_equals) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_is_equals) String() string {
	return "#<subr eq?>"
}

func (_ *_is_equals) IsList() bool {
	return false
}

func (i *_is_equals) Equals(sexp SExpression) bool {
	return i.TypeId() == sexp.TypeId()
}

func (_ *_is_equals) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength != 2 {
		return nil, errors.New("malformed if syntax")
	}

	first := args[0]

	second := args[1]

	return NewBool(first.Equals(second)), nil
}

func NewIsEq() SExpression {
	return &_is_equals{}
}
