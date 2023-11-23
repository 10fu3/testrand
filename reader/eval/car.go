package eval

import (
	"context"
	"errors"
)

//this implement is lisp car function

type _car struct{}

func (_ *_car) TypeId() string {
	return "subroutine.car"
}

func (_ *_car) AtomId() AtomType {
	return AtomTypeSubroutine
}

func (_ *_car) String() string {
	return "#<subr car>"
}

func (_ *_car) IsList() bool {
	return false
}

func (l *_car) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_car) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if 1 > argsLength {
		return nil, errors.New("need arguments size is 1")
	}

	if args[0].TypeId() != "cons_cell" {
		return nil, errors.New("need arguments type is list")
	}

	cons := args[0].(ConsCell)

	return cons.GetCar(), nil
}

func NewCar() SExpression {
	return &_car{}
}
