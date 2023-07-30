package eval

import "errors"

type _is_not struct {
}

func (_ *_is_not) Type() string {
	return "subroutine.is_equals"
}

func (_ *_is_not) String() string {
	return "#<subr eq?>"
}

func (_ *_is_not) IsList() bool {
	return false
}

func (i *_is_not) Equals(sexp SExpression) bool {
	return i.Type() == sexp.Type()
}

func (_ *_is_not) Apply(env Environment, arguments SExpression) (SExpression, error) {
	if "cons_cell" != arguments.Type() {
		return nil, errors.New("type error")
	}
	argCell := arguments.(ConsCell)

	first := argCell.GetCar()

	if "bool" != first.Type() {
		return first, nil
	}

	return NewBool(!first.(Bool).GetValue()), nil
}

func NewIsNot() SExpression {
	return &_is_not{}
}
