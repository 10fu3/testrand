package eval

import (
	"context"
	"errors"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func _syntax_transaction_Apply(self *Sexpression, ctx context.Context, env *Sexpression, args *Sexpression) (*Sexpression, error) {

	if SexpressionTypeConsCell != args._sexp_type_id {
		return CreateNil(), errors.New("type error")
	}

	cell := args._cell

	if !IsEmptyList(cell._cdr) {
		return CreateNil(), errors.New("need less than 2 params")
	}

	ok, err := env._env_superGlobalEnv.Transaction(func(stm concurrency.STM) error {
		_, err := Eval(context.WithValue(ctx, "transaction", stm), cell._car, env)
		if err != nil {
			return err
		}
		return nil
	})

	if !ok {
		return CreateNil(), err
	}

	return CreateEmptyList(), nil
}

func NewTransaction() *Sexpression {
	return CreateSpecialForm("transaction", _syntax_transaction_Apply)
}
