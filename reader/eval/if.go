package eval

import (
	"context"
	"errors"
	"fmt"
)

func _syntax_if_Apply(self *Sexpression, ctx context.Context, env *Sexpression, arguments *Sexpression) (*Sexpression, error) {
	args, argsLen, err := ToArray(arguments)

	if err != nil {
		return CreateNil(), err
	}

	if argsLen <= 1 || argsLen >= 4 {
		return CreateNil(), errors.New(fmt.Sprintf("too many argument: %d", len(args)))
	}
	argsIndex := 0

	statement := args[argsIndex]
	argsIndex++
	onTrue := args[argsIndex]

	evaluated, err := Eval(ctx, statement, env)
	if err != nil {
		return Eval(ctx, statement, env)
	}

	if evaluated.Equals(CreateBool(false)) {
		if argsLen == 2 {
			return CreateNil(), nil
		}
		argsIndex++
		onFalse := args[argsIndex]
		return Eval(ctx, onFalse, env)
	}
	return Eval(ctx, onTrue, env)
}

func NewIf() *Sexpression {
	return CreateSpecialForm("if", _syntax_if_Apply)
}
