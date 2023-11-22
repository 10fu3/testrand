package eval

import (
	"context"
	"errors"
)

func _subr_to_string_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {

	arr, arrSize, err := ToArray(args)

	if err != nil {
		return CreateNil(), err
	}

	if arrSize != 1 {
		return CreateNil(), errors.New("need arguments size is 1")
	}

	return CreateString(arr[0].String()), nil
}

func NewToString() *Sexpression {
	return CreateSubroutine("to-string", _subr_to_string_Apply)
}
