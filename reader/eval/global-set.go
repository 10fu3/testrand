package eval

import (
	"context"
	"errors"
	"fmt"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type _global_set struct {
}

func (_ *_global_set) TypeId() string {
	return "special_form.global_set"
}

func (_ *_global_set) AtomId() SExpressionType {
	return SExpressionTypeSpecialForm
}

func (_ *_global_set) String() string {
	return "#<syntax global_set>"
}

func (_ *_global_set) IsList() bool {
	return false
}

func (l *_global_set) Equals(sexp SExpression) bool {
	return l.TypeId() == sexp.TypeId()
}

func (_ *_global_set) Apply(ctx context.Context, env Environment, args []SExpression, argsLength uint64) (SExpression, error) {

	if argsLength != 2 {
		return nil, errors.New("need less than 2 params")
	}

	if args[0].TypeId() != "symbol" {
		return nil, errors.New("need 1st arguments type is symbol")
	}

	name := args[0].(Symbol)

	initValue := args[1].(ConsCell)

	if !IsEmptyList(initValue.GetCdr()) {
		return nil, errors.New("need less than 3 params")
	}
	evaluatedInitValue, err := Eval(ctx, initValue.GetCar(), env)

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	err = func() error {
		var err error
		if ctx.Value("transaction") != nil {
			transaction := ctx.Value("transaction").(concurrency.STM)
			if err != nil {
				return err
			}
			transaction.Put(fmt.Sprintf("/env/%s/%s", env.GetParentId(), name.String()), evaluatedInitValue.String())
		} else {
			err = env.GetSuperGlobalEnv().Put(fmt.Sprintf("/env/%s/%s", env.GetParentId(), name.String()), evaluatedInitValue.String(), nil)
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
