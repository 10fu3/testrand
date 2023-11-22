package eval

import (
	"context"
	"errors"
)

//this implement is lisp car function

func _sub_cdr_Apply(self *Sexpression, ctx context.Context, _ *Sexpression, args *Sexpression) (*Sexpression, error) {
	arr, arrSize, err := ToArray(args)
	if err != nil {
		return CreateNil(), err
	}
	if arrSize != 1 {
		return CreateNil(), errors.New("malformed car")
	}

	if SexpressionTypeConsCell != arr[0]._sexp_type_id {
		return CreateNil(), errors.New("malformed car")
	}

	cons := arr[0]._cell

	return cons._car, nil
}

func NewCdr() *Sexpression {
	return CreateSubroutine("cdr", _sub_cdr_Apply)
}
