package eval

import (
	"context"
	"errors"
)

type _symbol_name struct{}

func (s *_symbol_name) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {
	arr, err := ToArray(args)
	if err != nil {
		return nil, err
	}

	if 1 != len(arr) {
		return nil, errors.New("need arguments size is 1")
	}

	if "symbol" != arr[0].Type() {
		return nil, errors.New("need arguments type is symbol, but got " + arr[0].Type())
	}

	return NewString(arr[0].(Symbol).GetValue()), nil
}

func (s *_symbol_name) Type() string {
	return "subroutine.symbol_name"
}

func (s *_symbol_name) String() string {
	return "#<subr symbol_name>"
}

func (s *_symbol_name) IsList() bool {
	return false
}

func (s *_symbol_name) Equals(sexp SExpression) bool {
	return s.Type() == sexp.Type()
}

func NewSymbolName() SExpression {
	return &_symbol_name{}
}
