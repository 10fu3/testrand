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

func (_ *_string_split) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {
	arr, err := ToArray(args)
	if err != nil {
		return nil, err
	}
	if len(arr) < 1 {
		return nil, errors.New("need args size is 1")
	}

	if arr[0].TypeId() != "string" {
		return nil, errors.New("need args type is string")
	}

	str := arr[0].(Str).GetValue()

	var sep string
	if len(arr) == 2 {
		if arr[1].TypeId() != "string" {
			return nil, errors.New("need args type is string")
		}
		sep = arr[1].(Str).GetValue()
	} else {
		sep = ""
	}

	split := strings.Split(str, sep)

	rootList := &_cons_cell{}
	look := rootList

	if len(split) == 0 {
		return rootList, nil
	}

	for i := 0; i < len(split); i++ {
		look.Car = NewString(split[i])
		look.Cdr = NewConsCell(NewNil(), NewNil())
		if i+1 == len(split) {
			break
		}
		look = look.Cdr.(*_cons_cell)
	}

	return rootList, nil
}

func NewStringSplit() SExpression {
	return &_string_split{}
}
