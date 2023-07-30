package eval

import (
	"errors"
	"fmt"
)

type _loop struct{}

func (_ *_loop) Type() string {
	return "special_form.loop"
}

func (_ *_loop) String() string {
	return "#<syntax loop>"
}

func (_ *_loop) IsList() bool {
	return false
}

func (l *_loop) Equals(sexp SExpression) bool {
	return l.Type() == sexp.Type()
}

func (_ *_loop) Apply(env Environment, args SExpression) (SExpression, error) {
	if "cons_cell" != args.Type() {
		return nil, errors.New("need arguments")
	}
	arguments := args.(ConsCell)
	if !("cons_cell" == arguments.GetCdr().Type()) {
		return nil, errors.New("need arguments")
	}
	rawForms := arguments.GetCdr().(ConsCell)

	if !IsEmptyList(rawForms.GetCdr()) {
		return nil, errors.New("argument size must be 2, but got more than 2 arguments")
	}

	var evaluatedCond SExpression
	var err error

	for {
		evaluatedCond, err = Eval(arguments.GetCar(), env)
		if err != nil {
			return nil, err
		}

		if evaluatedCond.Type() != "bool" {
			return nil, errors.New("need 1st argument must be bool but got " + evaluatedCond.Type())
		}

		if !evaluatedCond.(Bool).GetValue() {
			return nil, nil
		}

		result, err := Eval(rawForms.GetCar(), env)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if result != nil {
			fmt.Println(result)
		}
	}
}

func NewLoop() SExpression {
	return &_loop{}
}
