package eval

import (
	"context"
	"errors"
	"strings"
)

type _string_split struct{}

func (_ *_string_split) TypeId() string {
	return "subroutine.string-split"
}

func (_ *_string_split) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSubroutine
}

func (_ *_string_split) String() string {
	return "#<subr string-split>"
}

func (_ *_string_split) IsList() bool {
	return false
}

func (s *_string_split) Equals(sexp SExpression) bool {
	return s.TypeId() == sexp.TypeId()
}

func (_ *_string_split) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {
	if argsLength < 1 {
		return nil, errors.New("need args size is 1")
	}

	if args[0].TypeId() != "string" {
		return nil, errors.New("need args type is string")
	}

	str := args[0].(Str).GetValue()

	var sep string
	if argsLength == 2 {
		if args[1].TypeId() != "string" {
			return nil, errors.New("need args type is string")
		}
		sep = args[1].(Str).GetValue()
	} else {
		sep = ""
	}

	split := strings.Split(str, sep)
	converted := make([]SExpression, len(split))
	for i := 0; i < len(split); i++ {
		converted[i] = NewString(split[i])
	}

	return &_native_array{Arr: converted}, nil
}

func NewStringSplit() SExpression {
	return &_string_split{}
}
