package eval

import (
	"context"
	"errors"
)

type _cons struct{}

func (_ *_cons) TypeId() string {
	return "subroutine.cons"
}

func (_ *_cons) AtomId() AtomType {
	return AtomTypeSubroutine
}

func (_ *_cons) String() string {
	return "#<subr cons>"
}

func (_ *_cons) IsList() bool {
	return false
}

func (_ *_cons) Equals(sexp SExpression) bool {
	return false
}

func (_ *_cons) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {
	if argsLength != 2 {
		return nil, errors.New("malformed cons")
	}

	return NewConsCell(args[0], args[1]), nil
}

func NewCons() SExpression {
	return &_cons{}
}
