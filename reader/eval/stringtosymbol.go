package eval

import (
	"context"
	"errors"
)

type _string_to_symbol struct{}

func (s *_string_to_symbol) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if 1 != argsLength {
		return nil, errors.New("need arguments size is 1")
	}

	if "string" != args[0].TypeId() {
		return nil, errors.New("need arguments type is string, but got " + args[0].TypeId())
	}

	return NewSymbol(args[0].(Str).GetValue()), nil
}

func (s *_string_to_symbol) TypeId() string {
	return "subroutine.string_to_symbol"
}

func (s *_string_to_symbol) AtomId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (s *_string_to_symbol) String() string {
	return "#<subr string_to_symbol>"
}

func (s *_string_to_symbol) IsList() bool {
	return false
}

func (s *_string_to_symbol) Equals(sexp SExpression) bool {
	return s.TypeId() == sexp.TypeId()
}

func NewStringToSymbol() SExpression {
	return &_string_to_symbol{}
}
