package eval

import (
	"errors"
	"fmt"
)

type _if struct {
}

func NewIf() SExpression {
	return &_if{}
}

func (_ *_if) Type() string {
	return "special_form.if"
}

func (_ *_if) String() string {
	return "#<syntax #if>"
}

func (_ *_if) IsList() bool {
	return false
}

func (i *_if) Equals(sexp SExpression) bool {
	return i.Type() == sexp.Type()
}

func (_ *_if) Apply(env Environment, argument SExpression) (SExpression, error) {
	args, err := ToArray(argument)

	if err != nil {
		return nil, err
	}

	if len(args) <= 1 || len(args) >= 4 {
		return nil, errors.New(fmt.Sprintf("too many argument: %d", len(args)))
	}
	argsIndex := 0

	statement := args[argsIndex]
	argsIndex++
	onTrue := args[argsIndex]

	evaluated, err := Eval(statement, env)
	if err != nil {
		return Eval(statement, env)
	}

	if evaluated.Equals(NewBool(false)) {
		if len(args) == 2 {
			return NewNil(), nil
		}
		argsIndex++
		onFalse := args[argsIndex]
		return Eval(onFalse, env)
	}
	return Eval(onTrue, env)
}
