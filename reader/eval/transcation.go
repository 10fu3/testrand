package eval

import (
	"context"
	"errors"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type _transaction struct{}

func (_ *_transaction) TypeId() string {
	return "special_form.transaction"
}

func (_ *_transaction) SExpressionTypeId() SExpressionType {
	return SExpressionTypeSpecialForm
}

func (_ *_transaction) String() string {
	return "#<syntax #transaction>"
}

func (_ *_transaction) IsList() bool {
	return false
}

func (t *_transaction) Equals(sexp SExpression) bool {
	return t.TypeId() == sexp.TypeId()
}

func (_ *_transaction) Apply(ctx context.Context, env Environment, args SExpression) (SExpression, error) {

	if "cons_cell" != args.TypeId() {
		return nil, errors.New("type error")
	}

	cell := args.(ConsCell)

	if !IsEmptyList(cell.GetCdr()) {
		return nil, errors.New("need less than 2 params")
	}

	ok, err := env.GetSuperGlobalEnv().Transaction(func(stm concurrency.STM) error {
		_, err := Eval(context.WithValue(ctx, "transaction", stm), cell.GetCar(), env)
		if err != nil {
			return err
		}
		return nil
	})

	if !ok {
		return nil, err
	}

	return NewConsCell(NewNil(), NewNil()), nil
}

func NewTransaction() SExpression {
	return &_transaction{}
}
