package eval

import (
	"context"
	"errors"
	"fmt"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func _syntax_global_set_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {
	if !args.IsConsCell() {
		return CreateNil(), errors.New("type error")
	}

	cell := args._cell

	if !cell._car.IsSymbol() {
		return CreateNil(), errors.New("need 1st arguments type is symbol")
	}

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

	if err != nil {
		return CreateNil(), err
	}
	err = func() error {
		if ctx.Value("transaction") != nil {
			transaction := ctx.Value("transaction").(concurrency.STM)
			transaction.Put(fmt.Sprintf("/env/%s/%s", env._env_parentId, name._string), evaluatedInitValue.String())
		} else {
			err = env._env_superGlobalEnv.Put(fmt.Sprintf("/env/%s/%s", env._env_parentId, name.String()), evaluatedInitValue.String(), nil)
		}
		return err
	}()

	if err != nil {
		return CreateNil(), err
	}
	return name, nil
}

func NewGlobalSet() *Sexpression {
	return CreateSpecialForm("global-set", _syntax_global_set_Apply)
}
