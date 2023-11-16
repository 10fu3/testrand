package eval

import (
	"context"
	"errors"
	"fmt"
)

type _loop struct{}

func (_ *_loop) TypeId() string {
	return "special_form.loop"
}

func (_ *_loop) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSpecialForm
}

func (_ *_loop) String() string {
	return "#<syntax loop>"
}

func (_ *_loop) IsList() bool {
	return false
}

func (l *_loop) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_loop) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {
	if "cons_cell" != args.TypeId() {
		return nil, errors.New("need arguments")
	}
	arguments := args.(ConsCell)
	if !("cons_cell" == arguments.GetCdr().TypeId()) {
		return nil, errors.New("need arguments")
	}
	rawForms := arguments.GetCdr().(ConsCell)

	if !IsEmptyList(rawForms.GetCdr()) {
		return nil, errors.New("argument size must be 2, but got more than 2 arguments")
	}

	var evaluatedCond SExpression
	var err error

	for {
		evaluatedCond, err = Eval(ctx, arguments.GetCar(), env)
		if err != nil {
			return nil, err
		}

		if evaluatedCond.TypeId() != "bool" {
			return nil, errors.New("need 1st argument must be bool but got " + evaluatedCond.TypeId())
		}

		if !evaluatedCond.(Bool).GetValue() {
			return nil, nil
		}

		_, err := Eval(ctx, rawForms.GetCar(), env)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}
}

func NewLoop() SExpression {
	return &_loop{}
}
