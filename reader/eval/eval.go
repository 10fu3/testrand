package eval

import (
	"errors"
	"strings"
)

func Eval(sexp SExpression, env Environment) (SExpression, error) {
	switch sexp.Type() {
	case "number":
		return sexp, nil
	case "bool":
		return sexp, nil
	case "nil":
		return sexp, nil
	case "symbol":
		return env.GetValue(sexp.(Symbol))
	case "cons_cell":
		cell := sexp.(ConsCell)
		applied, err := Eval(cell.GetCar(), env)
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(applied.Type(), "closure") || strings.HasPrefix(applied.Type(), "subroutine.") {
			appliedArgs, err := evalArgument(cell.GetCdr(), env)
			if err != nil {
				return nil, err
			}
			return applied.(Callable).Apply(env, appliedArgs)
		}
		if strings.HasPrefix(applied.Type(), "special_form.") {
			if err != nil {
				return nil, err
			}
			return applied.(Callable).Apply(env, cell.GetCdr())
		}

	}
	return nil, errors.New("unknown eval")
}

func evalArgument(sexp SExpression, env Environment) (SExpression, error) {
	if "cons_cell" != sexp.Type() {
		return Eval(sexp, env)
	}

	cell := sexp.(ConsCell)

	carEvaluated, err := Eval(cell.GetCar(), env)
	if err != nil {
		return nil, err
	}

	cdrEvaluated, err := evalArgument(cell.GetCdr(), env)
	if err != nil {
		return nil, err
	}

	return NewConsCell(carEvaluated, cdrEvaluated), nil
}
