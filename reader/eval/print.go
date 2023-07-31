package eval

import (
	"errors"
	"fmt"
)

type _println struct{}

func (_ *_println) Type() string {
	return "subroutine.print"
}

func (_ *_println) String() string {
	return "#<subr print>"
}

func (_ *_println) IsList() bool {
	return false
}

func (p *_println) Equals(sexp SExpression) bool {
	return p.Type() == sexp.Type()
}

func (_ *_println) Apply(env Environment, args SExpression) (SExpression, error) {
	arr, err := ToArray(args)
	if err != nil {
		return nil, err
	}
	if len(arr) != 1 {
		return nil, errors.New("need args size is 1")
	}
	fmt.Print(arr[0])
	return NewNil(), nil
}

func NewPrintln() SExpression {
	return &_println{}
}
