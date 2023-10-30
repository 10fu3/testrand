package eval

import (
	"context"
	"errors"
	"unicode/utf8"
)

type _string_len struct{}

func (_ *_string_len) Type() string {
	return "subroutine.string-len"
}

func (_ *_string_len) String() string {
	return "#<subr string-len>"
}

func (_ *_string_len) IsList() bool {
	return false
}

func (s *_string_len) Equals(sexp SExpression) bool {
	return s.Type() == sexp.Type()
}

func (_ *_string_len) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {
	arr, err := ToArray(args)
	if err != nil {
		return nil, err
	}
	if len(arr) != 1 {
		return nil, errors.New("need args size is 1")
	}

	size := int64(utf8.RuneCountInString(arr[0].(Str).GetValue()))

	return NewInt(size), nil
}

func NewStringLen() SExpression {
	return &_string_len{}
}
