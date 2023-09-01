package eval

import (
	"context"
	"errors"
)

type _define struct{}

func (_ *_define) Type() string {
	return "special_form.define"
}

func (_ *_define) String() string {
	return "#<syntax #define>"
}

func (_ *_define) IsList() bool {
	return false
}

func (d *_define) Equals(sexp SExpression) bool {
	return d.Type() == sexp.Type()
}

func onSymbolCall(ctx context.Context, env Environment, arguments SExpression) (SExpression, error) {

	if "cons_cell" != arguments.Type() {
		return nil, errors.New("type error")
	}

	cell := arguments.(ConsCell)

	name := cell.GetCar().(Symbol)

	if IsEmptyList(cell.GetCdr()) {
		env.Define(name, NewNil())
		return name, nil
	}

	initValue := cell.GetCdr().(ConsCell)

	if !IsEmptyList(initValue.GetCdr()) {
		return nil, errors.New("need less than 3 params")
	}
	evaluatedInitValue, err := Eval(ctx, initValue.GetCar(), env)

	if err != nil {
		return nil, err
	}
	env.Define(name, evaluatedInitValue)
	return name, nil
}

func (_ *_define) Apply(ctx context.Context, env Environment, arguments SExpression) (SExpression, error) {
	return onSymbolCall(ctx, env, arguments)
}

func NewDefine() SExpression {
	return &_define{}
}
