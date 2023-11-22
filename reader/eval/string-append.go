package eval

import (
	"context"
	"errors"
)

func _subr_string_append_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	arr, arrSize, err := ToArray(args)
	if err != nil {
		return CreateNil(), err
	}
	if arrSize < 1 {
		return CreateNil(), errors.New("need args size is 1")
	}

	var str string
	for i := uint64(0); i < arrSize; i++ {
		if SexpressionTypeString != arr[i]._sexp_type_id {
			return CreateNil(), errors.New("need args type is string, but got " + arr[i]._sexp_type_id.String())
		}
		str += arr[i]._string
	}
	return CreateString(str), nil
}

func NewStringAppend() *Sexpression {
	return CreateSubroutine("string-append", _subr_string_append_Apply)
}
