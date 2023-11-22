package eval

import (
	"context"
	"errors"
)

func _syntax__lambda_Apply(self *Sexpression, ctx context.Context, env *Sexpression, arguments *Sexpression) (*Sexpression, error) {
	args, argsLen, err := ToArray(arguments)
	if err != nil {
		return CreateNil(), err
	}

	if 2 != argsLen {
		return CreateNil(), errors.New("need arguments size is 2")
	}

	params := args[0]
	body := args[1]

	formalsArr, formalsSize, err := ToArray(params)

	if err != nil {
		return CreateNil(), err
	}

	closure, err := CreateClosure(body, formalsArr, env, formalsSize)

	if err != nil {
		return CreateNil(), err
	}

	return closure, nil
}

func NewLambda() *Sexpression {
	return CreateSpecialForm("lambda", _syntax__lambda_Apply)
}
