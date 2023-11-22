package eval

import (
	"context"
	"errors"
)

type _symbol_name struct{}

func (s *_symbol_name) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if 1 != argsLength {
		return nil, errors.New("need arguments size is 1")
	}

	if "symbol" != args[0].TypeId() {
		return nil, errors.New("need arguments type is symbol, but got " + args[0].TypeId())
	}

	return NewString(args[0].(Symbol).GetValue()), nil
}

func (s *_symbol_name) TypeId() string {
	return "subroutine.symbol_name"
}

func (s *_symbol_name) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (s *_symbol_name) String() string {
	return "#<subr symbol_name>"
}

func (s *_symbol_name) IsList() bool {
	return false
}

func (s *_symbol_name) Equals(sexp SExpression) bool {
	return s.TypeId() == sexp.TypeId()
}

func NewSymbolName() SExpression {
	return &_symbol_name{}
}
