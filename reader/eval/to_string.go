package eval

import (
	"context"
	"errors"
)

type _to_string struct{}

func (s *_to_string) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength != 1 {
		return nil, errors.New("need arguments size is 1")
	}

	return NewString(args[0].String()), nil
}

func (s *_to_string) TypeId() string {
	return "subroutine.to_string"
}

func (s *_to_string) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (s *_to_string) String() string {
	return "#<subr to_string>"
}

func (s *_to_string) IsList() bool {
	return false
}

func (s *_to_string) Equals(sexp SExpression) bool {
	return s.TypeId() == sexp.TypeId()
}

func NewToString() SExpression {
	return &_to_string{}
}
