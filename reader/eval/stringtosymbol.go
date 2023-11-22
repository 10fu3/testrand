package eval

import (
	"context"
	"errors"
)

func _subr_string_to_symbol_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	arr, arrSize, err := ToArray(args)
	if err != nil {
		return CreateNil(), err
	}

	if 1 != arrSize {
		return CreateNil(), errors.New("need arguments size is 1")
	}

	if SexpressionTypeString != arr[0]._sexp_type_id {
		return CreateNil(), errors.New("need arguments type is string, but got " + arr[0]._sexp_type_id.String())
	}

	return CreateSymbol(arr[0]._string), nil
}

func NewStringToSymbol() *Sexpression {
	return CreateSubroutine("string->symbol", _subr_string_to_symbol_Apply)
}
