package eval

import (
	"context"
	"errors"
)

func _syntax_set_Apply(self *Sexpression, ctx context.Context, env *Sexpression, arguments *Sexpression) (*Sexpression, error) {
	if SexpressionTypeConsCell != arguments._sexp_type_id {
		return CreateNil(), errors.New("type error")
	}

	cell := arguments._cell

	name := cell._car._symbol

	if IsEmptyList(cell._cdr) {
		return CreateNil(), errors.New("need 3rd arguments")
	}

	initValue := cell._cdr._cell

	if !IsEmptyList(initValue._cdr) {
		return CreateNil(), errors.New("need less than 3 params")
	}
	evaluatedInitValue, err := Eval(ctx, initValue._car, env)

	if err != nil {
		return CreateNil(), err
	}

	env._env_frame.Set(name._string, evaluatedInitValue)
	if err != nil {
		return CreateNil(), err
	}
	return name, nil
}

func NewSet() *Sexpression {
	return CreateSpecialForm("set!", _syntax_set_Apply)
}
