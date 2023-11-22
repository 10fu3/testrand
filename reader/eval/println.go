package eval

import (
	"context"
	"errors"
	"fmt"
)

type _println struct{}

func (_ *_println) TypeId() string {
	return "subroutine.print"
}

func (_ *_println) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_println) String() string {
	return "#<subr print>"
}

func (_ *_println) IsList() bool {
	return false
}

func (p *_println) Equals(sexp SExpression) bool {
	return p.TypeId() == sexp.TypeId()
}

func (_ *_println) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {
	if argsLength < 1 {
		return nil, errors.New("need args size is 1")
	}

	for i := uint64(0); i < argsLength; i++ {
		fmt.Print(args[i])
		if i+1 == argsLength {
			fmt.Println()
			break
		}
		fmt.Print(" ")
	}
	return NewNil(), nil
}

func NewPrintln() SExpression {
	return &_println{}
}
