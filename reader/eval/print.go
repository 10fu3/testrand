package eval

import (
	"context"
	"errors"
	"fmt"
)

type _print struct{}

func (_ *_print) TypeId() string {
	return "subroutine.print"
}

func (_ *_print) AtomId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_print) String() string {
	return "#<subr print>"
}

func (_ *_print) IsList() bool {
	return false
}

func (p *_print) Equals(sexp SExpression) bool {
	return p.TypeId() == sexp.TypeId()
}

func (_ *_print) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {
	if argsLength != 1 {
		return nil, errors.New("need args size is 1")
	}
	fmt.Print(args[0])
	return NewNil(), nil
}

func NewPrint() SExpression {
	return &_print{}
}
