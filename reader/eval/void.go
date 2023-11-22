package eval

import (
	"context"
	"errors"
)

func _syntax_void_Apply(self *Sexpression, ctx context.Context, env *Sexpression, arguments *Sexpression) (*Sexpression, error) {
	args, argsLen, err := ToArray(arguments)

	if err != nil {
		return CreateNil(), err
	}

	if 1 != argsLen {
		return CreateNil(), errors.New("need arguments size is 1")
	}

	_, err = Eval(ctx, args[0], env)

	if err != nil {
		return CreateNil(), err
	}

	return CreateNil(), nil
}

func NewVoid() *Sexpression {
	return CreateSpecialForm("void", _syntax_void_Apply)
}
