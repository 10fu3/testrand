package eval

import (
	"context"
	"errors"
	"fmt"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type _global_set struct {
}

func (_ *_global_set) Type() string {
	return "special_form.global_set"
}

func (_ *_global_set) String() string {
	return "#<syntax global_set>"
}

func (_ *_global_set) IsList() bool {
	return false
}

func (l *_global_set) Equals(sexp SExpression) bool {
	return l.Type() == sexp.Type()
}

func (_ *_global_set) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {
	if "cons_cell" != args.Type() {
		return nil, errors.New("type error")
	}

	cell := args.(ConsCell)

	name := cell.GetCar().(Symbol)

	if IsEmptyList(cell.GetCdr()) {
		return nil, errors.New("need 3rd arguments")
	}

	initValue := cell.GetCdr().(ConsCell)

	if !IsEmptyList(initValue.GetCdr()) {
		return nil, errors.New("need less than 3 params")
	}
	evaluatedInitValue, err := Eval(ctx, initValue.GetCar(), env)

	if err != nil {
		return nil, err
	}

	err = func() error {
		var err error
		if ctx.Value("transaction") != nil {
			transaction := ctx.Value("transaction").(concurrency.STM)
			transaction.Put(fmt.Sprintf("/env/%s/%s", env.GetId(), name.String()), evaluatedInitValue.String())
		} else {
			err = env.GetSuperGlobalEnv().Put(fmt.Sprintf("/env/%s/%s", env.GetId(), name.String()), evaluatedInitValue.String())
		}
		return err
	}()

	if err != nil {
		return nil, err
	}
	return name, nil
}

func NewGlobalSet() SExpression {
	return &_global_set{}
}
