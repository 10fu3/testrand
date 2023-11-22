package eval

import (
	"context"
	"errors"
)

func _syntax_define_Apply(self *Sexpression, ctx context.Context, env *Sexpression, arguments *Sexpression) (*Sexpression, error) {
	arr, arrSize, err := ToArray(arguments)

	if err != nil {
		return CreateNil(), err
	}

	if arrSize != 2 {
		return CreateNil(), errors.New("malformed define")
	}

	name := arr[0]
	initValue := arr[1]

	if SexpressionTypeSymbol != name.SexpressionTypeId() {
		return CreateNil(), errors.New("malformed define")
	}

	evaluatedInitValue, err := Eval(ctx, initValue, env)

	if err != nil {
		return CreateNil(), err
	}
	env._env_frame.Set(name._symbol._string, evaluatedInitValue)
	return name, nil
}

func NewDefine() *Sexpression {
	return CreateSpecialForm("define", _syntax_define_Apply)
}
