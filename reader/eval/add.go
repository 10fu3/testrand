package eval

import (
	"context"
	"errors"
)

type _add struct{}

func (_ *_add) Type() string {
	return "subroutine.add"
}

func (_ *_add) String() string {
	return "#<subr +>"
}

func (_ *_add) IsList() bool {
	return false
}

func (a *_add) Equals(sexp SExpression) bool {
	return a.Type() == sexp.Type()
}

func (_ *_add) Apply(ctx context.Context, _ Environment, args SExpression) (SExpression, error) {
	if "cons_cell" != args.Type() {
		return nil, errors.New("need arguments")
	}

	if IsEmptyList(args) {
		return NewInt(0), nil
	}

	arr, err := ToArray(args)

	if err != nil {
		return nil, err
	}
	arrIndex := 1

	var result Number = nil
	if "number" != arr[0].Type() {
		return nil, errors.New("need arguments type is number, but got " + arr[0].Type())
	}
	result = arr[0].(Number)
	for ; arrIndex < len(arr); arrIndex++ {
		if "number" != arr[arrIndex].Type() {
			return nil, errors.New("need arguments type is number, but got " + arr[arrIndex].Type())
		}
		result = NewInt(result.GetValue() + arr[arrIndex].(Number).GetValue())
	}
	return result, nil
}

func NewAdd() SExpression {
	return &_add{}
}
