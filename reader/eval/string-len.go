package eval

import (
	"context"
	"errors"
	"unicode/utf8"
)

type _string_len struct{}

func (_ *_string_len) TypeId() string {
	return "subroutine.string-len"
}

func (_ *_string_len) AtomId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_string_len) String() string {
	return "#<subr string-len>"
}

func (_ *_string_len) IsList() bool {
	return false
}

func (s *_string_len) Equals(sexp SExpression) bool {
	return s.TypeId() == sexp.TypeId()
}

func (_ *_string_len) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength != 1 {
		return nil, errors.New("need args size is 1")
	}

	size := int64(utf8.RuneCountInString(args[0].(Str).GetValue()))

	return NewInt(size), nil
}

func NewStringLen() SExpression {
	return &_string_len{}
}
