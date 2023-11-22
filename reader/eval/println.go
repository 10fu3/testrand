package eval

import (
	"context"
	"errors"
	"fmt"
)

func _subr_println_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	arr, arrSize, err := ToArray(args)
	if err != nil {
		return CreateNil(), err
	}
	if arrSize < 1 {
		return CreateNil(), errors.New("need args size is 1")
	}

	for i := uint64(0); i < arrSize; i++ {
		fmt.Print(arr[i])
		if i+1 == arrSize {
			fmt.Println()
			break
		}
		fmt.Print(" ")
	}
	return CreateNil(), nil
}

func NewPrintln() *Sexpression {
	return CreateSubroutine("println", _subr_println_Apply)
}
