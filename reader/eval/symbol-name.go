package eval

import (
	"context"
	"errors"
)

func _subr_symbol_name_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	arr, arrSize, err := ToArray(args)
	if err != nil {
		return CreateNil(), err
	}

	if 1 != arrSize {
		return CreateNil(), errors.New("need arguments size is 1")
	}

	if !arr[0].IsSymbol() {
		return CreateNil(), errors.New("need arguments type is symbol, but got " + arr[0]._sexp_type_id.String())
	}

	return CreateString(arr[0]._symbol._string), nil
}

func NewSymbolName() *Sexpression {
	return CreateSubroutine("symbol-name", _subr_symbol_name_Apply)
}
