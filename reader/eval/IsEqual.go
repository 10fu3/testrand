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

func (_ *_is_equals) Apply(ctx context.Context, env Environment, arguments SExpression) (SExpression, error) {
	if "cons_cell" != arguments.TypeId() {
		return nil, errors.New("type error")
	}
	argCell := arguments.(ConsCell)

	first := argCell.GetCar()

	if "cons_cell" != argCell.GetCdr().TypeId() {
		return nil, errors.New("type error")
	}

	second := argCell.GetCdr().(ConsCell)

	if !IsEmptyList(second.GetCdr()) {
		return nil, errors.New("argument size error")
	}

	return NewBool(first.Equals(second.GetCar())), nil
}

func NewIsEq() SExpression {
	return &_is_equals{}
}
