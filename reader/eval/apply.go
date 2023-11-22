package eval

import (
	"context"
	"errors"
)

func _sub_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	arr, arrSize, err := ToArray(args)
	if err != nil {
		return CreateNil(), err
	}
	if arrSize != 2 {
		return CreateNil(), errors.New("malformed apply")
	}
	car := arr[0]
	cdr := arr[1]

	if !car.IsCallable() {
		return CreateNil(), errors.New("first argument is not callable")
	}

	return car._applyFunc(self, ctx, env, cdr)
}

func NewApply() *Sexpression {
	return CreateSubroutine("apply", _sub_Apply)
}
