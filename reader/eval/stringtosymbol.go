package eval

import (
	"context"
	"errors"
)

type _string_to_symbol struct{}

func (s *_string_to_symbol) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {
	arr, err := ToArray(args)
	if err != nil {
		return nil, err
	}

	if 1 != len(arr) {
		return nil, errors.New("need arguments size is 1")
	}

	if "string" != arr[0].Type() {
		return nil, errors.New("need arguments type is string, but got " + arr[0].Type())
	}

	return NewSymbol(arr[0].(Str).GetValue()), nil
}

func (s *_string_to_symbol) Type() string {
	return "subroutine.string_to_symbol"
}

func (s *_string_to_symbol) String() string {
	return "#<subr string_to_symbol>"
}

func (s *_string_to_symbol) IsList() bool {
	return false
}

func (s *_string_to_symbol) Equals(sexp SExpression) bool {
	return s.Type() == sexp.Type()
}

func NewStringToSymbol() SExpression {
	return &_string_to_symbol{}
}
