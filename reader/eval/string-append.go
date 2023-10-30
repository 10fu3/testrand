package eval

import (
	"context"
	"errors"
)

type _string_append struct{}

func (_ *_string_append) Type() string {
	return "subroutine.string-append"
}

func (_ *_string_append) String() string {
	return "#<subr string-append>"
}

func (_ *_string_append) IsList() bool {
	return false
}

func (s *_string_append) Equals(sexp SExpression) bool {
	return s.Type() == sexp.Type()
}

func (_ *_string_append) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {
	arr, err := ToArray(args)
	if err != nil {
		return nil, err
	}
	if len(arr) < 1 {
		return nil, errors.New("need args size is 1")
	}

	var str string
	for i := 0; i < len(arr); i++ {
		str += arr[i].String()
	}
	return NewString(str), nil
}

func NewStringAppend() SExpression {
	return &_string_append{}
}
