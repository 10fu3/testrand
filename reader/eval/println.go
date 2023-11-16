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

func (_ *_println) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {
	arr, err := ToArray(args)
	if err != nil {
		return nil, err
	}
	if len(arr) < 1 {
		return nil, errors.New("need args size is 1")
	}

	for i := 0; i < len(arr); i++ {
		fmt.Print(arr[i])
		if i+1 == len(arr) {
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
