package eval

import (
	"context"
	"errors"
)

type _set struct{}

func (_ *_set) TypeId() string {
	return "special_form.set"
}

func (_ *_set) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSpecialForm
}

func (_ *_set) String() string {
	return "#<syntax #set>"
}

func (_ *_set) IsList() bool {
	return false
}

func (s *_set) Equals(sexp SExpression) bool {
	return s.TypeId() == sexp.TypeId()
}

func (_ *_set) Apply(ctx context.Context, env Environment, arguments SExpression) (SExpression, error) {
	if "cons_cell" != arguments.TypeId() {
		return nil, errors.New("type error")
	}

	cell := arguments.(ConsCell)

	name := cell.GetCar().(Symbol)

	if IsEmptyList(cell.GetCdr()) {
		return nil, errors.New("need 3rd arguments")
	}

	initValue := cell.GetCdr().(ConsCell)

	if !IsEmptyList(initValue.GetCdr()) {
		return nil, errors.New("need less than 3 params")
	}
	evaluatedInitValue, err := Eval(ctx, initValue.GetCar(), env)

	if err != nil {
		return nil, err
	}

	err = env.Set(name, evaluatedInitValue)
	if err != nil {
		return nil, err
	}
	return name, nil
}

func NewSet() SExpression {
	return &_set{}
}
