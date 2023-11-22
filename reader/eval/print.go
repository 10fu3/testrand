package eval

import (
	"context"
	"errors"
	"fmt"
)

func _subr_print_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	arr, arrSize, err := ToArray(args)
	if err != nil {
		return CreateNil(), err
	}
	if arrSize != 1 {
		return CreateNil(), errors.New("need args size is 1")
	}
	fmt.Print(arr[0])
	return CreateNil(), nil
}

func NewPrint() *Sexpression {
	return CreateSubroutine("print", _subr_print_Apply)
}
