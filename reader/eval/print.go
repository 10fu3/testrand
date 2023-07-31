package eval

import (
	"errors"
	"fmt"
)

type _print struct{}

func (_ *_print) Type() string {
	return "subroutine.print"
}

func (_ *_print) String() string {
	return "#<subr print>"
}

func (_ *_print) IsList() bool {
	return false
}

func (p *_print) Equals(sexp SExpression) bool {
	return p.Type() == sexp.Type()
}

func (_ *_print) Apply(env Environment, args SExpression) (SExpression, error) {
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

func NewPrint() SExpression {
	return &_println{}
}
